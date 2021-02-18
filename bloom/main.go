package main

import (
	"flag"
	"fmt"
	bf "github.com/willf/bloom"
	"hash/crc32"
	"strconv"

	//fnv "hash/fnv"
	//"runtime"
	//	xxhash "github.com/cespare/xxhash"
	"sync"
	"time"
)

var Falpro int = 0
var lock sync.Mutex

//arry chan size
const ArrySize = 5

var ArryBloomS [ArrySize]*bf.BloomFilter
var ArryChanData [ArrySize]chan []byte
var crcTable = crc32.MakeTable(crc32.IEEE)
var ArryChan_select [ArrySize]chan []byte

//
var cpunum int
var insertnum int
var selectnum int

func init() {
	flag.IntVar(&cpunum, "cpunum", 3, "cpu num parallel")
	flag.IntVar(&insertnum, "insertnum", 1e8, "cpu num parallel")
	flag.IntVar(&selectnum, "selectnum", 1e8, "cpu num parallel")

}

var exitchan chan interface{} = make(chan interface{}, cpunum)

func main() {

	flag.Parse()
	for i, _ := range ArryBloomS {
		ArryBloomS[i] = bf.New(600e8, 20)
	}

	for i, _ := range ArryChanData {
		ArryChanData[i] = make(chan []byte, 10000)
	}

	for i, _ := range ArryChan_select {
		ArryChan_select[i] = make(chan []byte, 10000)
	}
	//time.Sleep(time.Second*100)
	//var C chan []byte = make(chan []byte , 10000)

	t01 := time.Now()
	fmt.Println("cpunum:", cpunum)
	//for workid := 1; workid <= cpunum; workid++ {
	//	go Work_insert_bloom(workid, C, exitchan)
	//}

	for i := 0; i < ArrySize; i++ {
		go func(i int) {
			for {
				v, ok := <-ArryChanData[i]
				if !ok {
					break
				}
				ArryBloomS[i].Add(v)
			}
			exitchan <- struct{}{}
		}(i)

	}

	func() {
		for i := 0; i < insertnum; i++ {
			//	C <- []byte(strconv.Itoa(i))
			if i%1e7 == 0 {
				fmt.Printf("num:%v time:%v\n", i, time.Since(t01))
				t01 = time.Now()
			}

			tmp := strconv.Itoa(i)
			stobyte := []byte(tmp)
			hashVal := crc32.Checksum(stobyte, crcTable)
			index := int(hashVal) % 5
			ArryChanData[index] <- stobyte
		}
	}()

	for i := 0; i < ArrySize; i++ {
		close(ArryChanData[i])

	}

	for i := 0; i < cpunum; i++ {
		<-exitchan
	}

	fmt.Println(time.Since(t01))

	fmt.Println("select ...#################")

	T1 := time.Now()

	for w := 0; w < 5; w++ {
		go Work_select_bloom(w, ArryChan_select[w], exitchan)
	}

	for i := 0; i < selectnum; i++ {
		tmp := strconv.Itoa(i)
		//stobyte:=public.Str2bytes(tmp)
		stobyte := []byte(tmp)
		hashVal := crc32.Checksum(stobyte, crcTable)
		index := int(hashVal) % 5
		ArryChan_select[index] <- stobyte
	}

	for i := 0; i < ArrySize; i++ {
		close(ArryChan_select[i])
	}

	for i := 0; i < cpunum; i++ {
		<-exitchan
	}

	fmt.Println("falpro:", Falpro, time.Since(T1))

}
