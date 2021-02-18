package main

import (
	"fmt"
	"time"
)

func Work_select_bloom(workid int, jobchandata <-chan []byte, exitchan chan<- interface{}) {

	fmt.Printf("select workid:%v begin ...\n", workid)
	t0 := time.Now()
	num := 0
	for ii := range jobchandata {
		num++
		if num%10000000 == 0 {
			fmt.Printf("workid: %v num:%v  time:%v\n", workid, num, time.Since(t0))
			t0 = time.Now()

		}
		/*	tmp := strconv.Itoa(ii)
			hashVal := crc32.Checksum([]byte(tmp), crcTable)
			index := int(hashVal) % 5*/
		if ArryBloomS[workid].Test(ii) {

		} else {
			fmt.Println("tmp:::", string(ii))

			lock.Lock()
			Falpro++
			lock.Unlock()
		}

	}
	fmt.Printf("workid:%v end  \n", workid)
	exitchan <- struct{}{}
}
