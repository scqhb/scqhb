package public

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func ReadfileToChan(dirlist string, cc chan string) {
	dir, err := ioutil.ReadDir(dirlist)
	if err != nil {
		fmt.Println("read dir err:", err)
		panic(err)
	}
	var Ch8_filepath chan string = make(chan string, 10000)
	for _, file := range dir {
		if !file.IsDir() {
			filefullpath := fmt.Sprintf("%v/%v", dirlist, file.Name())
			Ch8_filepath <- filefullpath
			fmt.Println("list file:", filefullpath)
		}
	}
	close(Ch8_filepath)
	linecount := 0

lable001:
	for {
		select {
		case filepath, ok := <-Ch8_filepath:
			if !ok {
				break lable001
			}
			func() {
				t0 := time.Now()
				fi, err := os.Open(filepath)
				if err != nil {
					fmt.Println("读取文件错误:", filepath, err)
					panic(err)
				}
				defer fi.Close()
				br := bufio.NewReader(fi)
				for {
					linetmp, _, err := br.ReadLine()
					if err == io.EOF {
						break
					}

					linecount++

					if linecount%10000 == 0 {
						fmt.Println(linecount, time.Since(t0))
						t0 = time.Now()
					}

					cc <- string(linetmp)

				}
			}()
		default:
		}
	}
	fmt.Printf("file read com!")
	close(cc)
}
