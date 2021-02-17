package until

import (
	"bufio"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var Ctx = context.Background()

//Scxml is 制造xml文件
func Redis_Data(input_file string) {
	///home/oracle/gosrc/src/scdata/city.txt
	//cite文件加载进入redis

	rdb99 := redis.NewClient(&redis.Options{
		//	Network:            "",
		Addr:               "9.9.9.99:6379",
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           "",
		DB:                 10,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolSize:           10,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
		Limiter:            nil,
	})
	pong, err := rdb99.Ping(Ctx).Result()
	if err != nil {
		fmt.Println("redis ping errr", pong, err)
	}

	cityfile, err := os.Open(input_file)
	if err != nil {
		fmt.Println("err4:", err)
	}
	defer cityfile.Close()
	reader := bufio.NewReader(cityfile)
	bz := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if line != nil {
			bz++

			err2 := rdb99.Set(Ctx, strconv.Itoa(bz), string(line), 0).Err()
			if err != nil {
				fmt.Println("err2:", err2)
			}
		}

	}
	//名字文件加载进入redis
	/*namefile, err := os.Open("/home/oracle/gosrc/src/scdata/person_name.txt")
	if err!=nil{
		fmt.Println("err5:",err)
	}
	defer namefile.Close()
	bzname:=4000
	read2 := bufio.NewReader(namefile)
	for{
		line, _, err := read2.ReadLine()
		if err==io.EOF{
			break
		}else if err!=nil{
			fmt.Println("err7:",err)
			return
		}
		if line != nil {
			bzname++
			err2 := rdb99.Set(Ctx, strconv.Itoa(bzname), string(line), 0).Err()
			if err != nil {
				fmt.Println("err2:", err2)
			}
		}

	}*/
}

//生成xml文件
func Scxml(xmloutput string) {
	type Server struct {
		//	Text    string `xml:",chardata"`
		TRADEID string `xml:"trade_id"`
		XMLID   string `xml:"xml_id"`
		SEQNO   string `xml:"seq_no"`
		CSNM    string `xml:"csnm"`
		HTDT    string `xml:"htdt"`
		CRCD    string `xml:"crcd"`
		TICD    string `xml:"ticd"`
		REDT    string `xml:"redt"`
	}

	type Struct_xx struct {
		XMLName xml.Name `xml:"type"`
		Text    string   `xml:",chardata"`
		Version string   `xml:"version,attr"`
		Servers []Server `xml:"singleserver"`
	}

	rdb99 := redis.NewClient(&redis.Options{
		//	Network:            "",
		Addr:               "9.9.9.99:6379",
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           "",
		DB:                 10,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolSize:           10,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
		Limiter:            nil,
	})
	pong, err := rdb99.Ping(Ctx).Result()
	if err != nil {
		fmt.Println("redis ping errr", pong, err)
	}
	////v := &T_DHTR_TSDT_B{Version: "1111"}
	Wfile, err := os.OpenFile(xmloutput, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if err != nil {
		fmt.Println("err8:", err)
	}
	defer Wfile.Close()

	wbufio := bufio.NewWriter(Wfile)

	wbufio.WriteString(xml.Header)
	/////

	vvv := &Struct_xx{
		XMLName: xml.Name{},
		Text:    "",
		Version: "1111",
		Servers: nil,
	}

	start0 := time.Now()
	for i := 0; i < 200000; i++ {
		ss := strconv.Itoa(rand.Intn(3500) + 1)
		//	fmt.Println("ss:::",ss)
		result, err := rdb99.Get(Ctx, ss).Result()
		if err != nil {
			fmt.Println("err55:", err)
		}
		//fmt.Println("redis_get_value::", result)
		split := strings.Split(result, " ")
		//fmt.Println(split[0], split[1])
		sname := strconv.Itoa(rand.Intn(1000000) + 4001)
		ss_name, err := rdb99.Get(Ctx, sname).Result()
		if err != nil {
			fmt.Println("err60:", err, sname)
		}
		//生成500内的数据
		TICD := strconv.Itoa(rand.Intn(500))
		//生成A-Z的数字
		REDT2 := string(byte(rand.Intn(26) + 65))
		//初始化结构体Struct_T_DHTR_TSDT_B
		tmp := Server{Uuidv4(), Uuidv4(), strconv.Itoa(i), split[0], split[1], ss_name, TICD, REDT2}
		vvv.Servers = append(vvv.Servers, tmp)
		//v.Nodedatas = append(v.Nodedatas, Stmp)

	}
	output, err := xml.MarshalIndent(vvv, " ", " ")
	if err != nil {
		fmt.Printf("error7: %v\n", err)
	}

	wbufio.Write(output)
	wbufio.Flush()
	fmt.Println("fang fa 1", time.Since(start0))
	//////////////////##########
}

/*<type version="1111">
<singleserver>
*/

func ReadXml(inputxml string, outfiletxt string) {
	//定义表结构题
	type Struct_T_DHTR_TSDT_B struct {
		XMLName xml.Name `xml:"type"`
		Text    string   `xml:",chardata"`
		Version string   `xml:"version,attr"`
		Server  []struct {
			Text    string `xml:",chardata"`
			TRADEID string `xml:"trade_id"`
			XMLID   string `xml:"xml_id"`
			SEQNO   string `xml:"seq_no"`
			CSNM    string `xml:"csnm"`
			HTDT    string `xml:"htdt"`
			CRCD    string `xml:"crcd"`
			TICD    string `xml:"ticd"`
			REDT    string `xml:"redt"`
		} `xml:"singleserver"`
	}

	openfilexml, err := os.Open(inputxml)
	if err != nil {
		fmt.Println("err8:", err)
	}
	defer openfilexml.Close()
	all, err := ioutil.ReadAll(openfilexml)
	if err != nil {
		fmt.Println("err9:", err)
	}

	vv := Struct_T_DHTR_TSDT_B{}
	err10 := xml.Unmarshal(all, &vv)
	if err != nil {
		fmt.Println("err10:", err10)

	}
	//fmt.Println(vv)
	WxmlFile, err := os.OpenFile(outfiletxt, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("err111", err)
	}
	writerbuf := bufio.NewWriter(WxmlFile)

	for _, v1 := range vv.Server {
		out := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v|%v\n", v1.REDT, v1.CSNM, v1.CRCD, v1.HTDT, v1.CRCD, v1.TICD, v1.TRADEID, v1.XMLID)
		writerbuf.WriteString(out)
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
