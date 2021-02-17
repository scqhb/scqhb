package main

import (
	"fmt"
	"hash/crc32"
	"time"
)

func Work_insert_bloom(workid int, jobchandata <-chan []byte, exitchan chan<- interface{}) {

	fmt.Printf("workid:%v begin ...\n", workid)
	t0 := time.Now()
	num := 0
	for ii := range jobchandata {
		num++
		if num%10000000 == 0 {
			fmt.Printf("workid: %v num:%v  time:%v\n", workid, num, time.Since(t0))
			t0 = time.Now()
		}
  		hashVal := crc32.Checksum(ii, crcTable)
		index := int(hashVal) % 5
		ArryChanData[index] <- ii

		//	fmt.Printf("tmp:%v mm:%v\n",tmp,mm)

	}
	fmt.Printf("workid:%v end  \n", workid)
	exitchan <- struct{}{}
}
