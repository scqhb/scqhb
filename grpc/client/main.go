package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/gammazero/workerpool"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	proto "scdata/grpc" // 根据proto文件自动生成的代码
	"strings"
	"sync"
	"time"
)

/////////////////////////////

type (
	RoundRobinConf struct {
		Key       string
		lbkey     string
		Endpoints []string
	}
)

const (
	key   string = "grpcstreamserverkey"
	lbkey string = "lbkey002"
)

type (
	FileInfo struct {
		filename string
		filesize int
	}
)

//文件及路径缓存通道
var Ch_listfile chan FileInfo

//在目录中读取文件到chan
func SerchFileToChan(dirlist string) *chan FileInfo {
	chlist := make(chan FileInfo, 10000)

	var filecount int
	dir, err := ioutil.ReadDir(dirlist)
	if err != nil {
		fmt.Println("读取目录错误", err)
	}
	//读取目录文件
	for _, file := range dir {
		if !file.IsDir() {
			filefullpath := fmt.Sprintf("%v/%v", dirlist, file.Name())
			chlist <- FileInfo{
				filename: filefullpath,
				filesize: int(file.Size()) / 1024 / 1024,
			}

			filecount++

			fmt.Println("input file:", filefullpath)
		}
	}
	close(chlist)
	fmt.Println("文件目录读取完成,文件个数为:", filecount)
	return &chlist
}

var Ch_CacheLine chan string = make(chan string, 1000000)

//读取文件内容到chancache
func ReadFileToChanCache() {

	chan_file := SerchFileToChan(filedir)

	start0 := time.Now()
	linecount := 0
	func() {
	lablefor001:
		for {
			select {
			case fp, ok := <-*chan_file:
				if !ok {
					log.Printf("all file read compled!")
					break lablefor001
				}
				log.Printf("begin read file:%v file size:%v MB\n", fp.filename, fp.filesize)
				sfile, err := os.Open(fp.filename)

				if err != nil {
					log.Printf("read file err:%v %v", fp.filename, err)
					break
				}
				defer sfile.Close()
				br := bufio.NewReader(sfile)
				for {
					line, _, err := br.ReadLine()
					if err == io.EOF {
						log.Printf("read file compled:%v  size:%vMB", fp.filename, fp.filesize)
						break
					}
					linecount++
					if linecount%1000000 == 0 {
						fmt.Println("countline:", linecount, time.Since(start0))
						start0 = time.Now()
					}
					Ch_CacheLine <- string(line)
				}
			}
		}
		fmt.Printf("send count:%v  time:%v\n", linecount, time.Since(start0))
	}()
	close(Ch_CacheLine)
}

var lock sync.Mutex

func streamEtcdClient(etcds []string) proto.ChatClient {
	lock.Lock()
	defer lock.Unlock()
	var EClient proto.ChatClient
	//pb.SingleRpcBloomServiceClient
	fmt.Println("etcd server:", etcds)
	//连接etcd服务
	conf := &RoundRobinConf{
		Key:       key,
		lbkey:     lbkey,
		Endpoints: etcds,
	}
	client3, err := clientv3.New(clientv3.Config{Endpoints: conf.Endpoints, DialTimeout: time.Second * 5})
	if err != nil {
		panic(err)
	}
	//ticker := time.NewTicker(3 * time.Second)

	allServerList, err := client3.Get(context.TODO(), conf.Key, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	lgs, err := client3.Get(context.TODO(), conf.lbkey)
	if err != nil {
		if err == rpctypes.ErrKeyNotFound {
			_, err = client3.Put(context.TODO(), conf.lbkey, "1111")
			if err != nil {
				panic(err)
			}
		}
		panic(err)
	}

	serverMap := make(map[int]*grpc.ClientConn)

	for i, kv := range allServerList.Kvs {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		conn, _ := grpc.DialContext(ctx, string(kv.Value), grpc.WithInsecure())
		serverMap[i] = conn
	}

	l := len(serverMap)
	if l > 0 {
		for _, kv := range lgs.Kvs {
			versionId := int(kv.Version)

			fmt.Println("versionid::", versionId)
			index := versionId % l
			log.Printf("id: %v 获取到节点:%v", index, serverMap[index].Target())
			client := proto.NewChatClient(serverMap[index])

			EClient = client
			//加工数据

			//加工数据结束

		}
		_, err = client3.Put(context.TODO(), conf.lbkey, "0000")
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("there no rpc server,servermap is:", l)
	}

	return EClient
}

var usge string = `
exlample:
client -fc fc -etcd 9.9.9.64:2379  -filedir /u01/data`
var fc string
var etcd, filedir string
var etcds []string

func init() {

	flag.StringVar(&fc, "fc", "default", usge)
	flag.StringVar(&etcd, "etcd", "9.9.9.64:2379", "etcd server ip ")
	flag.StringVar(&filedir, "filedir", "/u01/data", "file dir ")
}

func main() {
	flag.Parse()

	if fc == "fc" {
		etcds = strings.Split(etcd, ",")

	} else {
		fmt.Println(usge)

		return
	}

	//读取文件到chan
	go ReadFileToChanCache()

	var sendnum int
	var lock sync.Mutex

	dumlicatecount := 0
	wp := workerpool.New(10)

	for i := 0; i < 10; i++ {

		fmt.Println("i:", i)
		wp.Submit(func() {
			procnum := i

			fmt.Println("pronum:", procnum)
			client := streamEtcdClient(etcds)
			resp, err2 := client.BloomStream(context.TODO())
			if err2 != nil {
				log.Println("errr2:", err2)
			}
			//加工数据
			go func() {
				tmp := 0
			lable001:
				for {
					select {
					case ctmp, ok := <-Ch_CacheLine:
						if !ok {
							break lable001
						}
						lock.Lock()
						sendnum++
						tmp++

						/*			xx:=&proto.Request{
											Input:                ctmp,

										}
									//	datamashal, err := proto2.Marshal(xx)
						*/
						lock.Unlock()

						err2 := resp.Send(&proto.Request{Input: ctmp})
						if err2 != nil {
							fmt.Println(err2)
						}
						if sendnum%10000000 == 0 {
							fmt.Println("客户端关闭流发送")
							resp.CloseSend()
							log.Printf("进程:%v  处理数据:%v", procnum, tmp)
							return
						}
					}
				}
			}()

			//加工数据结束

			for {
				fmt.Println("接收server 端消息")
				////	//respon 接收从 服务端返回的数据流
				_, err := resp.Recv()
				if err == io.EOF {
					log.Println("收到服务端的结束信号")
					break //如果收到结束信号，则退出“接收循环”，结束客户端程序
				}
				if err != nil {
					log.Println("接收数据出错:", err)
				}

				lock.Lock()
				dumlicatecount++
				lock.Unlock()
				// 没有错误的情况下，打印来自服务端的消息
				//	fmt.Printf(respon.Output)
				log.Printf("  dumlicatecount111:%v sendnum:%v  fasl: %2.8f\n", dumlicatecount, sendnum, float64(dumlicatecount)/float64(sendnum))

			}
			log.Printf("  dumlicatecount222:%v sendnum:%v  fasl: %2.8f\n", dumlicatecount, sendnum, float64(dumlicatecount)/float64(sendnum))

		})
		time.Sleep(time.Second * 1)
	}

	wp.StopWait()

}

//http restful 客户端
func main_httprestful() {

	de := http.DefaultTransport
	transport, ok := de.(*http.Transport)
	if !ok {
		fmt.Println("transport xxxx")
		return
	}
	transport.MaxIdleConns = 1000
	transport.MaxIdleConnsPerHost = 1000
	SerchFileToChan("/u01/data")

	start0 := time.Now()
	falsenum := 0
	linecount := 0
	func() {
	lablefor001:
		for {
			select {
			case fp, ok := <-Ch_listfile:
				if !ok {
					break lablefor001
				}
				fmt.Println("读取文件:", fp.filename)
				sfile, err := os.Open(fp.filename)
				if err != nil {
					fmt.Println("读取文件错误", err)
				}
				defer sfile.Close()
				br := bufio.NewReader(sfile)
				for {
					line, _, err := br.ReadLine()
					if err == io.EOF {
						break
					}
					linecount++
					if linecount%10000 == 0 {
						fmt.Println("countline:", linecount, time.Since(start0))
						start0 = time.Now()
					}
					resp, err := http.Get(fmt.Sprintf("http://9.9.9.80:9999/Getfalsepositive/%v", string(line)))
					//time.Sleep(time.Second)
					if err != nil {
						fmt.Println("http get err", err)
						return
					}
					io.Copy(ioutil.Discard, resp.Body)
					if resp.StatusCode == 400 {
						falsenum++

						fmt.Println("falsenum:", falsenum)

					}
					resp.Body.Close()

				}
			}
		}
		fmt.Printf("send count:%v  time:%v\n", linecount, time.Since(start0))
	}()

}
