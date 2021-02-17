package until

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type RBIF struct { //报告基本信息
	RINM string `xml:"RINM"` //报告机构名称
	RICD string `xml:"RICD"` //报告机构编码
	FIRC string `xml:"FIRC"` //报告机构所在地区编码
	LEI  string `xml:"LEI"`  //LEI编码
	CTTN string `xml:"CTTN"` //交易客户总数
}

type CBIF struct { //客户基本信息

	CSNM  string    `xml:"CSNM"` //客户号                 #主键
	CTVC  string    `xml:"CTVC"` //客户职业（行业）类别
	CRNM  string    `xml:"CRNM"` //法定代表人姓名（对公客户）
	CRIT  string    `xml:"CRIT"` //法定代表人身份证件类型（对公客户）
	ORIT  string    `xml:"ORIT"` //法定代表人其他身份证件类型（对公客户）
	CRID  string    `xml:"CRID"` //法定代表人身份证件号码（对公客户）
	CCIF  CCIF_CBIF `xml:"CCIF"`
	CTNTs CTNTs     `xml:"CTNTs"`
}

type CATI struct { //客户和交易信息
	Seqno string `xml:"seqno,attr"`
	CBIF  CBIF   `xml:"CBIF"` //客户基本信息
	HTDT  string `xml:"HTDT"` //大额交易发生日期                                       #主键
	HTCRs HTCRs  `xml:"HTCRs"`
}

type HTCR struct {
	Seqno string `xml:"seqno,attr"`
	CRCD  string `xml:"CRCD"` //大额交易特征代码                                     #主键
	TTNM  string `xml:"TTNM"` //交易总数
	CCIFs CCIFs  `xml:"CCIFs"`
}

type ATIF struct {
	CBAT string `xml:"CBAT"` //客户银行账户类型
	CBAC string `xml:"CBAC"` //客户银行账号
	CABM string `xml:"CABM"` //客户银行账号开户银行名称
	CTAT string `xml:"CTAT"` //客户支付账户类型
	CTAC string `xml:"CTAC"` //客户支付账号////////////
	OATM string `xml:"OATM"` //客户开户时间
	CPIN string `xml:"CPIN"` //客户所在非银行支付机构名称
	CPBA string `xml:"CPBA"` //客户所在非银行支付机构银行账号
	CPBN string `xml:"CPBN"` //客户所在非银行支付机构银行账号开户行名称
}

type TBIF struct {
	TBNM string `xml:"TBNM"` //代办人姓名
	TBIT string `xml:"TBIT"` //代办人身份证件(证明文件)类型
	OITP string `xml:"OITP"` //代办人其他身份证件(证明文件)类型
	TBID string `xml:"TBID"` //代办人身份证件(证明文件)号码
	TBNT string `xml:"TBNT"` //代办人国籍
}

type TSIF struct {
	TSTM string `xml:"TSTM"` //<TSTM>交易时间</TSTM>
	TRCD string `xml:"TRCD"` //<TRCD>交易发生地</TRCD>
	TICD string `xml:"TICD"` //<TICD>业务标识号</TICD>                  #主键
	CTTP string `xml:"CTTP"` //<CTTP>货币资金转移方式</CTTP>
	TSCT string `xml:"TSCT"` //<TSCT>涉外收支交易分类与代码</TSCT>
	TSDR string `xml:"TSDR"` //<TSDR>资金收付标志</TSDR>
	CRPP string `xml:"CRPP"` //<CRPP>资金用途</CRPP>
	CRTP string `xml:"CRTP"` //<CRTP>交易币种</CRTP>
	CRAT string `xml:"CRAT"` //<CRAT>交易金额</CRAT>
	CRMB string `xml:"CRMB"` //<CRMB>交易金额（折人民币）</CRMB>
	CUSD string `xml:"CUSD"` //<CUSD>交易金额（折美元）</CUSD>
	TMNM string `xml:"TMNM"` //<TMNM>交易商品名称</TMNM>
	OCTT string `xml:"OCTT"` //<OCTT>非柜台交易方式</OCTT>
	OCEC string `xml:"OCEC"` //<OCEC>非柜台交易方式的设备代码</OCEC>
	OOCT string `xml:"OOCT"` //<OOCT>其他非柜台交易方式</OOCT>
	BPTC string `xml:"BPTC"` //<BPTC>银行与非银行支付机构之间的业务交易编码</BPTC>
	PMTC string `xml:"PMTC"` //<PMTC>非银行支付机构与商户之间的业务交易编码</PMTC>
	CTIP string `xml:"CTIP"` //<CTIP>客户交易IP地址</CTIP>
}

type TCIF struct { //交易对手基本信息
	TCNM string `xml:"TCNM"` //<TCNM>交易对手姓名（名称）</TCNM>
	TCIT string `xml:"TCIT"` //<TCIT>交易对手身份证件（证明文件）类型</TCIT>
	OITP string `xml:"OITP"` //<OITP>交易对手其他身份证件（证明文件）类型</OITP>
	TCID string `xml:"TCID"` //<TCID>交易对手身份证件（证明文件）号码</TCID>
	TCAT string `xml:"TCAT"` //<TCAT>交易对手银行账户类型</TCAT>
	TCBA string `xml:"TCBA"` //<TCBA>交易对手银行账号</TCBA>
	TCBN string `xml:"TCBN"` //<TCBN>交易对手银行账号开户银行名称</TCBN>
	TCTT string `xml:"TCTT"` //<TCTT>交易对手支付账户类型</TCTT>
	TCTA string `xml:"TCTA"` //<TCTA>交易对手支付账号</TCTA>
	TCPN string `xml:"TCPN"` //<TCPN>交易对手所在非银行支付机构名称</TCPN>
	TCPA string `xml:"TCPA"` //<TCPA>交易对手所在非银行支付机构银行账号</TCPA>
	TPBN string `xml:"TPBN"` //<TPBN>交易对手所在非银行支付机构银行账号开户银行名称</TPBN>
	TCIP string `xml:"TCIP"` //<TCIP>交易对手的交易IP地址</TCIP>
}

type TSDT struct {
	Seqno string `xml:"seqno,attr"`
	FICD  string `xml:"FICD"` //报告机构分支机构（网点）代码
	ATIF  ATIF   `xml:"ATIF"`
	TBIF  TBIF   `xml:"TBIF"`
	TSIF  TSIF   `xml:"TSIF"`
	TCIF  TCIF   `xml:"TCIF"`
	ROTFs ROTFs  `xml:"ROTFs"`
}
type ROTF struct {
	Seqno string `xml:"seqno,attr"`
	Text  string `xml:",chardata"`
}
type ROTFs struct {
	ROTF []ROTF `xml:"ROTF"`
}

type CCIF struct {
	Seqno string `xml:"seqno,attr"`
	CTNM  string `xml:"CTNM"` //客户名称（姓名）
	CITP  string `xml:"CITP"` //客户身份证件（证明文件）类型
	OITP  string `xml:"OITP"` //客户其他身份证件（证明文件）类型
	CTID  string `xml:"CTID"` //客户身份证件（证明文件）号码
	TSDTs TSDTs  `xml:"TSDTs"`
}
type CCIFs struct { //客户身份信息
	CCIF []CCIF `xml:"CCIF"`
}

type TSDTs struct { //大额交易信息
	TSDT []TSDT `xml:"TSDT"`
}

////////###########

type CCTL struct { //客户联系电话
	Text  string `xml:",chardata"`
	Seqno string `xml:"seqno,attr"` //客户联系电话
}

type CTAR struct { //客户住址/经营地址
	Text  string `xml:",chardata"`
	Seqno string `xml:"seqno,attr"`
}

type CCEI struct { //客户其他联系方式
	Text  string `xml:",chardata"`
	Seqno string `xml:"seqno,attr"` //客户其他联系方式
}

///////###########

type CCIF_CBIF struct { //客户身份信息
	CCTLs struct { //客户联系电话
		CCTL []CCTL `xml:"CCTL"`
	} `xml:"CCTLs"`
	CTARs struct { //客户住址/经营地址
		CTAR []CTAR `xml:"CTAR"`
	} `xml:"CTARs"`
	CCEIs struct { //客户其他联系方式
		CCEI []CCEI `xml:"CCEI"`
	} `xml:"CCEIs"`
}
type CTNT struct { //客户国籍
	Text  string `xml:",chardata"`
	Seqno string `xml:"seqno,attr"`
}
type CTNTs struct {
	CTNT []CTNT `xml:"CTNT"`
}

type HTCRs struct {
	HTCR []HTCR `xml:"HTCR"`
}

type CATIs struct {
	CATI []CATI `xml:"CATI"`
}

type PHTR2 struct {
	XMLName xml.Name `xml:"PHTR"`
	Text    string   `xml:",chardata"`
	RBIF    RBIF     `xml:"RBIF"`
	CATIs   CATIs    `xml:"CATIs"`
}

var MapBankAddr map[int]string = make(map[int]string)
var MapPepleAddr map[int]string = make(map[int]string)
var MapCountry map[int]string = make(map[int]string)
var MapPersonName map[int]string = make(map[int]string)
var MapWord map[int]string = make(map[int]string)
var MapGs map[int]string = make(map[int]string)

var MapErrinfo map[int]string = make(map[int]string)

//##########################xxxx
func Read_toMap(file string, m1 map[int]string) {
	oile, err := os.Open(file)
	if err != nil {
		fmt.Println("err12:", err)
	}
	reader := bufio.NewReader(oile)
	//m1:=make(map[int]string)
	banknum := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF || err != nil {
			break
		}

		m1[banknum] = string(line)
		banknum++

	}

}

func GetIpaddr() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

func Readfile_toMap() {

	MapErrinfo[0] = "异常扣款交易 0x001"
	MapErrinfo[1] = "异常转账交易 0x002"

	Read_toMap("asset/address.txt", MapPepleAddr)
	Read_toMap("asset/bank_addr.txt", MapBankAddr)
	Read_toMap("asset/country.txt", MapCountry)
	Read_toMap("asset/person_name.txt", MapPersonName)
	Read_toMap("asset/word.txt", MapWord)
	Read_toMap("asset/gs.txt", MapGs)

}
