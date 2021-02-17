package until

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

func Struct_to_file(phtr2 *PHTR2) {

	output, err := xml.MarshalIndent(phtr2, "", "  ")
	if err != nil {
		fmt.Printf("error7: %v\n", err)
	}
	rand.Seed(time.Now().UnixNano())
	fnum := rand.Intn(100000)
	filename := fmt.Sprintf("NPH%v-%v-%08d", phtr2.RBIF.RICD, time.Now().Format("20060102"), fnum)
	fullfilename := fmt.Sprintf("xml/%v.xml", filename)
	err2 := ioutil.WriteFile(fullfilename, output, 0600)
	if err2 != nil {
		fmt.Println("file write xml err:", err2)
	}
	Ch2_Xmlfile <- filename

}

func Marchl_Outfile(workid int, jobid chan int, result chan struct{}) {
	var Done chan struct{}
	func() {
	label:
		for {
			select {
			case p := <-Ch2_PHTR2:
				Struct_to_file(p)
				result <- struct{}{}
			case <-Done:
				break label
			}

		}

	}()

}
