package main

import (
	"encoding/json"
	"flag"
	"fmt"
	bf "github.com/willf/bloom"
	"io/ioutil"
	"runtime"
	grpc "scdata/grpc/server"
	"scdata/public"
	singlerpc "scdata/singlerpc/singleserver"
	"scdata/until"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var (
	fc, dirlist, single_ip_port, stream_ip_port string
	sctxtline                                   uint
	sctxtfiles, printnum, proc, port            int
	usge                                        string = `func options 
1.	scdata -fc sctxt -dirlist /data/dat1  -sctxtfiles 10 -sctxtline 1000000
2.	scdata -fc filetoredis -dirlist=/ul01/data  -printnum 100000  -proc 2
3.	scdata -fc filetoBloomFile -dirlist=/data/dat1  -printnum 100000  -proc 10 -stream_ip_port 9.9.9.80:8601 -single_ip_port 9.9.9.80:8501

filetoBloomFile
`
)

func init() {

	for i, _ := range public.ArryBloomS {
		public.ArryBloomS[i] = bf.New(100e8, 20)
	}

	for i, _ := range public.ArryChanData {
		public.ArryChanData[i] = make(chan []byte, 10000)
	}

	for i, _ := range public.ArryChan_select {
		public.ArryChan_select[i] = make(chan []byte, 10000)
	}

	flag.StringVar(&fc, "fc", "default", usge)
	flag.UintVar(&sctxtline, "sctxtline", 0, "sc txt file line size")
	flag.IntVar(&sctxtfiles, "sctxtfiles", 0, "sc txt file number")
	flag.IntVar(&printnum, "printnum", 0, "printnum number")
	flag.IntVar(&proc, "proc", 0, "proc number")
	flag.StringVar(&dirlist, "dirlist", ".", "dir path ")
	flag.IntVar(&port, "port", 8250, "server port")
	flag.StringVar(&single_ip_port, "single_ip_port", "9.9.9.80:8501", "single_ip_port: 9.9.9.80:8501")
	flag.StringVar(&stream_ip_port, "stream_ip_port", "9.9.9.80:8601", "stream_ip_port: 9.9.9.80:8601")
	for i := 0; i < 1e5; i++ {
		public.Map_jgbm[i] = fmt.Sprintf("F%v", string(until.Krandv2(13, 0)))
	}
}

func main() {

	flag.Parse()

	//生成txt文件
	if fc == "sctxt" && sctxtline > 100 && sctxtfiles > 0 {
		for numfile := 0; numfile < sctxtfiles; numfile++ {
			until.GetUuidToFile(dirlist, sctxtline)
		}
	}
	//将文本文件读取到redis的bloom里面
	if fc == "filetoredis" && dirlist != "" && printnum > 0 && proc > 0 {
		dir, err := ioutil.ReadDir(dirlist)
		if err != nil {
			fmt.Println("read dir err:", err)
			panic(err)
		}
		var Ch8_filepath chan string = make(chan string, 10000)
		for _, file := range dir {
			if !file.IsDir() {
				filefullpath := fmt.Sprintf("%v/%v", dirlist, file.Name())
				//filesize := file.Size() / 1024 / 1024
				Ch8_filepath <- filefullpath
				fmt.Println("input file:", filefullpath)
			}
		}
		close(Ch8_filepath)
		tt := time.Now()
		for i := 0; i < proc; i++ {
			until.Wg.Add(1)
			go until.TestBloomRedis(Ch8_filepath, printnum)
		}

		until.Wg.Wait()
		fmt.Println("总耗时: ", time.Since(tt))
	}
	//将文本文件读取到读取到内存bloom,并启动对外服务
	if fc == "filetoBloomFile" && dirlist != "" && printnum > 0 && proc > 0 {
		dir, err := ioutil.ReadDir(dirlist)
		if err != nil {
			fmt.Println("read dir err:", err)
			panic(err)
		}
		var Ch8_filepath chan string = make(chan string, 10000)

		for _, file := range dir {
			if !file.IsDir() {
				filefullpath := fmt.Sprintf("%v/%v", dirlist, file.Name())
				//filesize := file.Size() / 1024 / 1024
				Ch8_filepath <- filefullpath
				fmt.Println("input file:", filefullpath)
			}
		}
		close(Ch8_filepath)
		tt := time.Now()
		func() {
			for i := 0; i < public.ArrySize; i++ {
				go func(i int) {
					for {
						v, ok := <-public.ArryChanData[i]
						if !ok {
							break
						}
						public.ArryBloomS[i].Add(v)
					}
					public.Exitchan <- struct{}{}
				}(i)
			}
		}()

		for i := 0; i < proc; i++ {
			go until.ReadFileToBloom(i, Ch8_filepath, until.Filter)
		}

		for i := 0; i < public.ArrySize; i++ {
			<-public.Exitchan
		}
		for i := 0; i < public.ArrySize; i++ {
			close(public.ArryChanData[i])
		}

		fmt.Println("文本文件加载进内存bloom完成,总耗时: ", time.Since(tt))
		fmt.Println("开始对外提供服务")
		////go until.HttpServer()
		go singlerpc.SingleRpcMain(single_ip_port)
		grpc.BloomServiceGrpcStream(stream_ip_port)

		//	until.SerchFileFromBitmap("/u01/data/p20200911-1599797523126200076.txt", until.Filter)
	}

	/*
		runtime.GOMAXPROCS(runtime.NumCPU()/2)
			flag.Parse()
			start0 := time.Now()

			if *cpuprofile != "" {
				f, err := os.Create(*cpuprofile)
				if err != nil {
					fmt.Println("err create file fail", err)
				}

				pprof.StartCPUProfile(f)
				defer pprof.StopCPUProfile()

			}
			//KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
			//KC_RAND_KIND_UPPER = 2 // 大写字母
			//KC_RAND_KIND_NUM   = 0 // 纯数字
			until.Uncompresstxt()

			go func() {
				var Done chan struct{}
			lable:
				for {
					rand.Seed(time.Now().UnixNano())
					select {
					case until.Ch1_ALL <- string(until.Krand(24, until.KC_RAND_KIND_ALL)):
					case until.Ch1_UPPER <- string(until.Krand(24, until.KC_RAND_KIND_UPPER)):
					case until.Ch1_NUM <- string(until.Krand(24, until.KC_RAND_KIND_NUM)):
					case until.Ch1_NUM7 <- string(until.Krand(7, until.KC_RAND_KIND_NUM)):
					case until.Ch1_NUM4 <- string(until.Krand(4, until.KC_RAND_KIND_NUM)):
					case <-Done:
						break lable
					}
				}
				defer close(until.Ch1_ALL)
			}()

			until.Readfile_toMap()
			////

			var job chan int
			var result chan struct{} = make(chan struct{})

			for i := 0; i < 1; i++ {
				go until.Marchl_Outfile(i, job, result)
			}

			go func() {
				for i := 1; i < 1000; i++ {
					until.Sc(4, 3, 6000, 2)
					fmt.Println("生成结构体:", i)
				}

			}()
	*/

	/*	go func() {
			var Done chan struct{}
			var txt int =0
			label:
			for {
				select {
				case xmlname := <-until.Ch2_Xmlfile:
					until.Xml_jx(fmt.Sprintf("xml/%v.xml", xmlname), fmt.Sprintf("txt/%v.txt", xmlname))
					txt++
					fmt.Println("平面txt文件个数:",txt)
				case <-Done:
					break label
				}
			}

		}()

		for i := 1; i < 1000; i++ {
			<-result
			fmt.Println("生成xml文件个数:", i)
		}

		fmt.Println(time.Since(start0))
	*/
	/*Select_BitNum := uint(1e4)
	var otherwise, indata int64
	//定义bitmap将实际存储实际数据量:n

	arry2 := make([]string, 1)
	N2 := uint(1e7)
	//生成查询的数据记录时间
	t2 := time.Now()
	func() {
		for t2 := uint(0); t2 < Select_BitNum; t2++ {
			s2 := until.Uuidv4()
			arry2 = append(arry2, s2)
		}
	}()
	//查询数据生成结束时间
	fmt.Printf("生成查询数据 %d条,耗时%v \n", N2, time.Since(t2))

	//查bloom数据是否存在
	start3 := time.Now()
	for _, v := range arry2 {
		if Filter.Test([]byte(v)) {
			indata++
		} else {
			otherwise++
		}
	}

	fmt.Printf("检索查询数据 %d 条   耗时:%v\n", N2, time.Since(start3))

	fmt.Println("一定不存在:", otherwise)
	fmt.Println("可能存在:", indata-1)
	fmt.Printf("错误率: %10f    数据总数:%10d\n", float64(float64(indata-1)/float64(otherwise+indata-1)), otherwise+indata-1)

		fmt.Println("uint32 max::", int(^uint32(0)>>1))
		fmt.Println("uint64 max::", int(^uint64(0)>>1))
	*/
}

func main2() {

	/*		AddString("World").
			AddUInt16(uint16(16)).
			AddUInt32(uint32(32)).
			AddUInt64(uint64(64)).
			AddUint16Batch([]uint16{1, 2, 3})
		filter.Add([]byte("abcd"))*/
	/* 	fmt.Printf("Hello exist:%t\n", filter.Test([]byte("aaa")))
	   	fmt.Printf("World exist:%t\n", filter.TestString("bbb"))
	   	fmt.Printf("uint 16 exist:%t\n", filter.TestUInt16(uint16(16)))
	   	fmt.Printf("uint 32 exist:%t\n", filter.TestUInt32(uint32(32)))
	   	fmt.Printf("uint 64 exist:%t\n", filter.TestUInt64(uint64(64)))*/

}

/*
func main22(){
err:=until.ReadfiletoMap()
if err !=nil{
	fmt.Println("err:")
}


	XmlT3()
 }



func main2() {
	//	Test_XMLMarshal()
	XmlT3()
}

///



func XmlT3() {
	rand.Seed(time.Now().UnixNano())
	lenMapXm:=len(until.MapXm)
	lenMapCity:=len(until.MapCity)
	rand.Seed(time.Now().UnixNano())
	// until.MapXm[int64(rand.Intn(lenMapXm))]
	//until.MapCity[int64(rand.Intn(lenMapCity))]

	const (
		Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	)
	f, err := os.Create("myxml.xml")
	if err != nil {
		fmt.Println("err: ", err.Error())
		return
	}
	defer f.Close()
	f.Write([]byte(Header))
//是是是
	stu := Student{random(20,80), Address{until.MapCity[int64(rand.Intn(lenMapCity))], until.MapCity[int64(rand.Intn(lenMapCity))]}, []Email{Email{"home", "home@qq.com"}, Email{"work", "work@qq.com"}}, until.MapXm[int64(rand.Intn(lenMapXm))], "zhang"}
	fmt.Println("stu: ", stu)
	//
	encoder := xml.NewEncoder(f)
	err1 := encoder.Encode(stu)
	if err1 != nil {
		fmt.Println("err1: ", err1.Error())
		return
	}


	/////////////
	stu = Student{random(20,80), Address{until.MapCity[int64(rand.Intn(lenMapCity))], until.MapCity[int64(rand.Intn(lenMapCity))]}, []Email{Email{"home", "home@qq.com"}, Email{"work", "work@qq.com"}}, until.MapXm[int64(rand.Intn(lenMapXm))], "zhang"}
	fmt.Println("stu: ", stu)
	//
	encoder = xml.NewEncoder(f)
	err1 = encoder.Encode(stu)
	if err1 != nil {
		fmt.Println("err1: ", err1.Error())
		return
	}
	//解码xml
	//ÖØÖÃÎss
	f.Seek(0, os.SEEK_SET)
	decoder := xml.NewDecoder(f)
	var strName string
	for {
		token, err2 := decoder.Token()
		if err2 != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			stelm := xml.StartElement(t)
			fmt.Println("start: ", stelm.Name.Local)
			strName = stelm.Name.Local
		case xml.EndElement:
			endelm := xml.EndElement(t)
			fmt.Println("end: ", endelm.Name.Local)
		case xml.CharData:
			data := xml.CharData(t)
			str := string(data)
			switch strName {
			case "City":
				fmt.Println("city:", str)
			case "first":
				fmt.Println("first: ", str)
			}
		}
	}
}
*/

type Monitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	PauseTotalNs uint64

	NumGC        uint32
	NumGoroutine int
}

func NewMonitor(duration int) {
	var m Monitor
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)

		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		m.Alloc = rtm.Alloc
		m.TotalAlloc = rtm.TotalAlloc
		m.Sys = rtm.Sys
		m.Mallocs = rtm.Mallocs
		m.Frees = rtm.Frees

		// Live objects = Mallocs - Frees
		m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		m.PauseTotalNs = rtm.PauseTotalNs
		m.NumGC = rtm.NumGC

		// Just encode to json and print
		b, _ := json.Marshal(m)
		fmt.Println(string(b))
	}
}
