package until

import (
	"bufio"
	"fmt"
	redisbloom "github.com/RedisBloom/redisbloom-go"
	redis "github.com/gomodule/redigo/redis"
	"io"
	"os"
	"time"
)

func File_readis(filepath string, printnum int) {
	Pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "redishost:6379", redis.DialPassword(""))
		},
		TestOnBorrow:    nil,
		MaxIdle:         100,
		MaxActive:       100,
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}

	BloomClient001 := redisbloom.NewClientFromPool(Pool, "bloomclient001")
	err := BloomClient001.Reserve("largebloom", 0.000001, 1e10)
	if err != nil {
		fmt.Println("create largebloom fail:::", err)
	}
	start := time.Now()

	start0 := time.Now()
	fi, err := os.Open(filepath)
	if err != nil {
		fmt.Println("err1:", err)
		panic(err)
	}
	defer fi.Close()

	var count = 0
	var arraystring []string
	br := bufio.NewReader(fi)

	timecount := 0
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		count++
		if count%printnum == 0 {
			fmt.Printf("读取数据到bit: %d 条  放入bloom耗时:%v\n", count, time.Since(start0))
			start0 = time.Now()
		}

		arraystring = append(arraystring, string(line))
		timecount++
		if timecount == 500 {
			_, err3 := BloomClient001.BfAddMulti("largebloom", arraystring)
			if err3 != nil {
				fmt.Println("err33:", err3)
			}
			arraystring = arraystring[0:0]
			timecount = 0
		}
	}
	fmt.Printf("读取文件%v完成,耗时:%v\n", filepath, time.Since(start))
	fmt.Println("file jia gong wan cheng ")

}

//testredis sss
func TestBloomRedis(filepath chan string, printnum int) {
	defer Wg.Done()

lablefor001:
	for {
		select {
		case fp, ok := <-filepath:
			if !ok {
				break lablefor001
			}
			File_readis(fp, printnum)
		}
	}

}
