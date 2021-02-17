package main

import (
	"bufio"
	pb "scdata/singlerpc"
	"sync"

	"context"
	"fmt"
	"github.com/gammazero/workerpool"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"os"
	//	pb "scdata/singlerpc"
	"time"
)

/////###

type (
	RoundRobinConf struct {
		Key       string
		lbkey     string
		Endpoints []string
	}
	CreateRoundIndexFunc func() (interface{}, error)
	RoundRobin           struct {
		index   int
		lock    sync.Mutex
		targets []interface{}
	}
)

const (
	key   string = "grpcserverkey"
	lbkey string = "lbkey001"
)

////###

var Ch_Search_file chan string = make(chan string, 10000)

func SerchFileFromBitmap(dirlist string) {

	dir, err := ioutil.ReadDir(dirlist)
	if err != nil {
		fmt.Println("读取目录错误", err)
	}
	for _, file := range dir {
		if !file.IsDir() {
			filefullpath := fmt.Sprintf("%v/%v", dirlist, file.Name())
			//filesize := file.Size() / 1024 / 1024
			Ch_Search_file <- filefullpath
			fmt.Println("input file:", filefullpath)
		}
	}
	close(Ch_Search_file)
	fmt.Println("文件读取完成,文件个数为")
}

var Ch_CacheLine chan string = make(chan string, 1000000)

func ReadFileToChanCache() {

	SerchFileFromBitmap("/data")
	start0 := time.Now()
	linecount := 0
	func() {
	lablefor001:
		for {
			select {
			case fp, ok := <-Ch_Search_file:
				if !ok {
					break lablefor001
				}
				fmt.Println("读取文件:", fp)
				sfile, err := os.Open(fp)
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
					if linecount%1000000 == 0 {
						fmt.Println("countline:", linecount, time.Since(start0))
						start0 = time.Now()
					}
					Ch_CacheLine <- string(line)
					/*					resp, err2 := client.SingleRpc(context.Background(), &pb.ClientRequest{Requestmess: string(line)})
										if err2 != nil {
											log.Fatalf("singleclient resuest err:%v", err2)
											return
										}
										log.Printf("resp:%s", resp.Responsemess)*/

				}
			}
		}
		fmt.Printf("send count:%v  time:%v\n", linecount, time.Since(start0))
	}()
	close(Ch_CacheLine)
}

/*/////////////////////////////////
type SingleRpcBloomServiceClient interface {
	// Sends a greeting
	SingleRpc(ctx context.Context, in *ClientRequest, opts ...grpc.CallOption) (*ServerResponse, error)
}*/
var lock sync.Mutex

func SingleEtcdClient() pb.SingleRpcBloomServiceClient {
	lock.Lock()
	defer lock.Unlock()

	var EClient pb.SingleRpcBloomServiceClient

	//连接etcd服务
	conf := &RoundRobinConf{
		Key:       key,
		lbkey:     lbkey,
		Endpoints: []string{"9.9.9.64:2379"},
	}
	client3, err := clientv3.New(clientv3.Config{Endpoints: conf.Endpoints})
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
		ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
		conn, _ := grpc.DialContext(ctx, string(kv.Value), grpc.WithInsecure())
		serverMap[i] = conn
	}

	l := len(serverMap)
	if l > 0 {
		for _, kv := range lgs.Kvs {
			versionId := int(kv.Version)
			fmt.Println("versionid::", versionId)
			index := versionId % l
			fmt.Println("select node id:", index, serverMap[index])
			client := pb.NewSingleRpcBloomServiceClient(serverMap[index])

			if client != nil {
				fmt.Println(err)
			}
			EClient = client
			//加工数据

			//加工数据结束

		}
		fmt.Println("put............................................")
		_, err = client3.Put(context.TODO(), conf.lbkey, "0000")
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("there no rpc server,servermap is:", l)
	}

	return EClient
}

func main() {
	go ReadFileToChanCache()
	wp := workerpool.New(20)

	for i := 0; i < 20; i++ {
		fmt.Println("iii::::::::::::::::::::::::::", i)
		wp.Submit(func() {
			client := SingleEtcdClient()

			//加工数据
		lable001:
			for {
				select {
				case ctmp, ok := <-Ch_CacheLine:
					if !ok {
						break lable001
					}
					resp, err2 := client.SingleRpc(context.TODO(), &pb.ClientRequest{Requestmess: ctmp})
					if err2 != nil {
						log.Println(err2)
					}
					if resp.Responsemess == "true" {
						log.Println("success  ", resp.Responsemess)

					}
				}
			}
			//加工数据结束
		})
	}

	wp.StopWait()

}
