package until

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

func Xml_jx(inputxml string, outfiletxt string) {

	WxmlFile, err := os.OpenFile(outfiletxt, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("err111", err)
	}

	writerbuf := bufio.NewWriter(WxmlFile)

	//定义表结构题
	openfilexml, err := os.Open(inputxml)
	if err != nil {
		fmt.Println("err8:", err)
	}
	defer openfilexml.Close()
	all, err := ioutil.ReadAll(openfilexml)
	if err != nil {
		fmt.Println("err9:", err)
	}

	vv := &PHTR2{}
	err10 := xml.Unmarshal(all, vv)
	if err != nil {
		fmt.Println("err10:", err10)

	}

	for _, cati8 := range vv.CATIs.CATI {

		csnm := cati8.CBIF.CSNM
		htdt := cati8.HTDT
		for _, htcr8 := range cati8.HTCRs.HTCR {
			crcd := htcr8.CRCD

			for _, ccif := range htcr8.CCIFs.CCIF {
				for _, tsdt8 := range ccif.TSDTs.TSDT {
					ticd := tsdt8.TSIF.TICD
					out := fmt.Sprintf("%v,%v,%v,%v\n", csnm, htdt, crcd, ticd)
					writerbuf.WriteString(out)

				}

			}

		}

	}

	writerbuf.Flush()

	/*decoder := xml.NewDecoder(bytes.NewReader(readFile))
	for{
		t, err := decoder.Token()
		if err!=nil{
			if err==io.EOF{
				fmt.Printf("Parse xml finished!!!!!!!!1")
			}else{
				fmt.Printf("Failed to Parse xml %v\n",err)
			}
			break
		}
		t=xml.CopyToken(t)
		switch t:=t.(type) {
		case xml.StartElement:
			fmt.Printf("startElement:<%v>\n",t.Name.Local)
		case xml.EndElement:
			fmt.Printf("EndElement:<%v>\n",t.Name.Local)
		case xml.CharData:
			fmt.Printf("CharData:%v\n",string(t))
		case xml.Comment:
			fmt.Printf("comment:<!--%v-->\n",string(t))

		}
		fmt.Println("###########################################")



	}*/

}
