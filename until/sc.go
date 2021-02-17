package until

import (
	"encoding/xml"
	"fmt"
	//"io/ioutil"
	"math/rand"
	"scdata/asset"
	"strconv"
	"strings"
	"time"
)

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// 随机字符串
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

//释放文本文件到目录
func Uncompresstxt() {
	dirs := []string{"asset"} // 设置需要释放的目录
	for _, dir := range dirs {
		// 解压dir目录到当前目录
		if err := asset.RestoreAssets("./", dir); err != nil {
			fmt.Println("RestoreAssets ", dir, "error", "\n", err)
		}
	}

}

/////////////////////####################################
func Sc(CatiNum, HtcrNum, TsdtNum, CcifNum int) {
	/*	var CatiNum int = 3 //每份报告包含的客户基本信息数量
		var HtcrNum int = 4 //触发大额特征数量
		var TsdtNum int = 30000 //每个账户信息下包含的交易数量
		var CcifNum int = 2 //每个交易特征包含的客户身份数量
	*/

	//ATIF ################
	cbac := (<-Ch1_NUM)[:16]
	mapnum := rand.Intn(100000)
	cbat := strconv.Itoa(mapnum + 620000)
	num16 := (<-Ch1_NUM)[:16]

	cpba16 := string(Krand(16, KC_RAND_KIND_NUM))

	datetime := time.Now().Format("20060102150405")

	//RBIF 报告机构编码
	rand.Seed(time.Now().UnixNano())
	Large_letter := rand.Intn(26) + 'A'
	rcid := fmt.Sprintf("%c%014d", Large_letter, rand.Intn(1000000))
	rinm := fmt.Sprintf("报告机构%v", rcid)

	rbif := RBIF{
		RINM: rinm,
		RICD: rcid,
		FIRC: (<-Ch1_NUM)[:6],
		LEI:  (<-Ch1_ALL)[:20],
		CTTN: "2",
	}

	catis := CATIs{}
	for i := 1; i <= CatiNum; i++ {
		ccif_cbif := CCIF_CBIF{}
		for i := 0; i < 2; i++ {
			rand.Seed(time.Now().UnixNano())
			r := rand.Intn(96)
			ccei_addr := strings.Split(MapPepleAddr[r], "#")
			ctar_addr := strings.Split(MapPepleAddr[r+1], "#")
			cctl_addr := strings.Split(MapPepleAddr[r+2], "#")

			ccei := CCEI{
				Text:  ccei_addr[1],
				Seqno: ccei_addr[0],
			}
			ctar := CTAR{
				Text:  ctar_addr[1],
				Seqno: ctar_addr[0],
			}
			cctl := CCTL{
				Text:  cctl_addr[1],
				Seqno: cctl_addr[0],
			}
			ccif_cbif.CCEIs.CCEI = append(ccif_cbif.CCEIs.CCEI, ccei)
			ccif_cbif.CTARs.CTAR = append(ccif_cbif.CTARs.CTAR, ctar)
			ccif_cbif.CCTLs.CCTL = append(ccif_cbif.CCTLs.CCTL, cctl)

		}
		// ctnss
		ctnts := CTNTs{}
		for i := 0; i < 2; i++ {
			r := rand.Intn(229)
			ctnt := CTNT{
				Text:  MapCountry[r],
				Seqno: fmt.Sprintf("%03d", i),
			}
			ctnts.CTNT = append(ctnts.CTNT, ctnt)
		}

		rr := rand.Intn(9000)
		cbif := CBIF{
			CSNM:  <-Ch1_NUM7,
			CTVC:  "46656",
			CRNM:  MapPersonName[rr+3],
			CRIT:  "610001",
			ORIT:  "610001",
			CRID:  cpba16,
			CCIF:  ccif_cbif,
			CTNTs: ctnts,
		}

		htcrs := HTCRs{}
		for i := 1; i <= HtcrNum; i++ {

			ccifs := CCIFs{}
			for i := 1; i <= CcifNum; i++ {

				///存入TSDTS 5 个
				tsdts := TSDTs{}
				for i := 1; i <= TsdtNum; i++ {
					rand.Seed(time.Now().UnixNano())
					rr := rand.Intn(9000)
					rword := rand.Intn(120)
					mapnum := rand.Intn(250)
					rrgs := rand.Intn(390)
					//#############
					atif := ATIF{
						CBAT: cbat,
						CBAC: num16,
						CABM: strings.Split(MapBankAddr[mapnum], " ")[0],
						CTAT: cbat,
						CTAC: cbac,
						OATM: datetime,
						CPIN: MapGs[rrgs],
						CPBA: cpba16,
						CPBN: strings.Split(MapBankAddr[mapnum+1], " ")[0],
					}
					tbif := TBIF{
						TBNM: MapPersonName[rr],
						TBIT: "610001",
						OITP: "610001",
						TBID: num16,
						TBNT: "NOR",
					}
					tsif := TSIF{
						TSTM: datetime,
						TRCD: "CCK466561",
						TICD: <-Ch1_ALL,
						CTTP: "0203",
						TSCT: cbac,
						TSDR: "01",
						CRPP: MapWord[rword],
						CRTP: "GBP",
						CRAT: "4665.62",
						CRMB: "4665.62",
						CUSD: "706.911",
						TMNM: "@N",
						OCTT: "12",
						OCEC: "5QJPWT6MM6YKWV6N",
						OOCT: "12",
						BPTC: <-Ch1_ALL,
						PMTC: <-Ch1_ALL,
						CTIP: GetIpaddr(),
					}
					tcif := TCIF{
						TCNM: MapPersonName[rr+1],
						TCIT: "610001",
						OITP: "610001",
						TCID: num16,
						TCAT: "620001",
						TCBA: num16,
						TCBN: strings.Split(MapBankAddr[mapnum+2], " ")[0],
						TCTT: "620001",
						TCTA: num16,
						TCPN: MapGs[rrgs+1],
						TCPA: num16,
						TPBN: strings.Split(MapBankAddr[mapnum+3], " ")[0],
						TCIP: GetIpaddr(),
					}
					//#############

					rotfs := ROTFs{}
					for i := 0; i < 2; i++ {
						rotf := ROTF{fmt.Sprintf("%04d", i), MapErrinfo[i]}
						rotfs.ROTF = append(rotfs.ROTF, rotf)
					}

					tsdt := TSDT{
						Seqno: strconv.Itoa(i),
						FICD:  cbac,
						ATIF:  atif,
						TBIF:  tbif,
						TSIF:  tsif,
						TCIF:  tcif,
						ROTFs: rotfs,
					}
					tsdts.TSDT = append(tsdts.TSDT, tsdt)
				} ///////////////////

				ccif := CCIF{
					Seqno: "1",
					CTNM:  MapPersonName[rr+4],
					CITP:  "610001",
					OITP:  "610001",
					CTID:  "466561531580045320",
					TSDTs: tsdts,
				}
				ccifs.CCIF = append(ccifs.CCIF, ccif)

			}
			htcr := HTCR{
				Seqno: "1",
				CRCD:  <-Ch1_NUM4,
				TTNM:  "100",
				CCIFs: ccifs,
			}
			htcrs.HTCR = append(htcrs.HTCR, htcr)
		}

		cati := CATI{
			Seqno: strconv.Itoa(i),
			CBIF:  cbif,
			HTDT:  time.Now().Format("20060102150405"), //主键大额交易发生日期
			HTCRs: htcrs,
		}
		catis.CATI = append(catis.CATI, cati)

	}

	Ch2_PHTR2 <- &PHTR2{
		XMLName: xml.Name{},
		Text:    "",
		RBIF:    rbif,
		CATIs:   catis,
	}

	//////////////////##########

	//fmt.Println(string(output))

}
