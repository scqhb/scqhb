package until

import (
	"bufio"
	"scdata/public"

	//"context"
	"fmt"
	"os"
	//"sync"
	"time"
)

func ErrCheck(err error) {
	if err != nil {
		fmt.Println("err:::", err)
		return
	}
}

//调用uuid函数并将获取的数据写入文件
func GetUuidToFile(dirlist string, num uint) {
	var C chan *SS_fivecol = make(chan *SS_fivecol, 5000)
	go Rand_Line(C)
	//file name
	////filename := fmt.Sprintf("%v/p%v-%08d.txt", dirlist, time.Now().Format("20060102"), time.Now().UnixNano())
	tmp:=public.GetFilenum()
	filename := fmt.Sprintf("%v/testdb.bwallcol01.%04d.csv", dirlist,tmp)
	createfile, err := os.Create(filename)
	ErrCheck(err)
	defer createfile.Close()

	w := bufio.NewWriter(createfile) //创建新的 Writer 对象
	var Count uint
	t0 := time.Now()
	////_, err3 := w.WriteString(fmt.Sprintf("ricd\tocnm\totcd\totic\totdt\tprocid\n"))
	_, err3 := w.WriteString(fmt.Sprintf("ricd\n"))

	ErrCheck(err3)

lablefor:
	for {
		select {
		case tmpss := <-C:
			{
				////_, err3 := w.WriteString(fmt.Sprintf("%v%v%v%v%v%v\n", tmpss.Ch_bgjg, tmpss.Ch_khh, tmpss.Ch_tzm, tmpss.Ch_yswbh, tmpss.Ch_jyrq, tmpss.Ch_pronum))

				//_, err3 := w.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\n", tmpss.Ch_bgjg, tmpss.Ch_khh, tmpss.Ch_tzm, tmpss.Ch_yswbh, tmpss.Ch_jyrq, tmpss.Ch_pronum))
				_, err3 := w.WriteString(fmt.Sprintf("%v%v%v%v%v%v\n", tmpss.Ch_bgjg, tmpss.Ch_khh, tmpss.Ch_tzm, tmpss.Ch_yswbh, tmpss.Ch_jyrq, tmpss.Ch_pronum))

				ErrCheck(err3)
				Count++
				if Count%num == 0 {
					fmt.Printf("data line:%v time:%v \n", Count, time.Since(t0))
					t0 = time.Now()
					//cancelFunc()
					break lablefor
				}
			}

		default:
			//	fmt.Println("no found")
		}

	}

	defer w.Flush()
	//	Wg.Wait()

}
