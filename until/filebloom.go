package until

import (
	"bufio"
	"fmt"
	"github.com/willf/bloom"
	"hash/crc32"
	"io"
	"os"
	"scdata/public"
	"time"
)

//创建布隆过滤器,传入bitmap大小和hash个数
var N_BitArrySize = uint(1e2)

//定义m/n比值常数
var MNconst = uint(30)

//定义hash 函数个数:k
var K_HashNum = uint(20)

//定义查询数据数据条数

//并发加载数据量
var Parallel = uint(1)

//加载数据时间记时t0

var Filter = bloom.New(MNconst*N_BitArrySize, K_HashNum)

var MapXm map[int64]string = make(map[int64]string)
var MapCity map[int64]string = make(map[int64]string)

//按照行读取文件到bitmap
func ReadFileToBloom(workid int, filepath chan string, Fbit *bloom.BloomFilter) {
	fmt.Printf("workid:%v begin ...", workid)
	for {
		filein, ok := <-filepath
		if !ok {
			fmt.Printf("workid: %v Read file complete", workid)
			break
		}

		func() {
			count := 0
			start0 := time.Now()
			//N_BitArrySize, MNconst, K_HashNum, Parallel
			fmt.Printf("真实数据量:%d  数组长度/真实数据:%d  hash个数:%d    并发:%d\n", N_BitArrySize, MNconst, K_HashNum, Parallel)
			//开始装载数据
			//	arry1=append(arry1,s1)
			fi, err := os.Open(filein)
			if err != nil {
				fmt.Println("读取文件错误:", filein, err)
				panic(err)
			}
			defer fi.Close()
			br := bufio.NewReader(fi)
			for {
				line, _, err := br.ReadLine()
				if err == io.EOF {
					break
				}
				count++
				if count%1e7 == 0 {
					fmt.Printf("workid: %v 读取数据到bit: %d 条  放入bloom耗时:%v\n", workid, count, time.Since(start0))
					start0 = time.Now()
				}
				hashVal := crc32.Checksum(line, public.CrcTable)
				index := int(hashVal) % public.ArrySize
				public.ArryChanData[index] <- line
			}
			fmt.Println("读取file到bit完成t0", filein)
		}()
	}

	public.Exitchan <- struct{}{}
}

//查找文件中的字符是否在bitmap里面
func SerchFileFromBitmap(searchfile string, Fbit *bloom.BloomFilter) {
	start0 := time.Now()
	falsepositive := 0

	//开始装载数据
	//	arry1=append(arry1,s1)
	sfile, err := os.Open(searchfile)
	ErrCheck(err)
	defer sfile.Close()
	var count = 0
	br := bufio.NewReader(sfile)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		count++
		if count%1e6 == 0 {
			fmt.Printf("检索数据: %d 条  检索bloom耗时:%v\n", count, time.Since(start0))
			start0 = time.Now()
		}

		if Fbit.Test(line) {
			falsepositive++
		}

	}

	fmt.Printf("假阳性数据条数%v  假阳率%10f", falsepositive, float64(float64(falsepositive)/float64(count)))

}
