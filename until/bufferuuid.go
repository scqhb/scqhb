package until

import (
	"bytes"
	"github.com/google/uuid"
	fastrand "github.com/valyala/fastrand"
	"runtime"
	"scdata/public"

	"strconv"
	"time"
)

/*
type readOp int8

const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid   readOp = 0  // Non-read operation.
	opReadRune1 readOp = 1  // Read rune of size 1.
	opReadRune2 readOp = 2  // Read rune of size 2.
	opReadRune3 readOp = 3  // Read rune of size 3.
	opReadRune4 readOp = 4  // Read rune of size 4.
)

type Buffer struct {
	b  bytes.Buffer
	rw sync.RWMutex
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	b.rw.RUnlock()
	defer b.rw.RUnlock()
	return b.b.Read(p)
}
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.b.Write(p)
}

func (b *Buffer) WriteString(s string) (n int, err error) {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.b.WriteString(s)

}

func (b *Buffer) String() string {
	b.rw.RUnlock()
	defer b.rw.RUnlock()
	return b.b.String()
}

func BufferUuidv4() string {

	var buffer bytes.Buffer
	str, _ := uuid.NewRandom()
	b := str.String()
	buffer.WriteString(b)
	str, _ = uuid.NewRandom()
	b = str.String()
	buffer.WriteString(b)
	return buffer.String()
}
*/
//BufferUuidv4 生成uuid,根据num传入的值拼接字符串长度,num为1时候返回的字符串长度就是36个字符
func BufferUuidv4() string {
	var buffer8 bytes.Buffer
	str1, _ := uuid.NewRandom()
	str2, _ := uuid.NewRandom()
	buffer8.WriteString(strconv.Itoa(time.Now().Nanosecond()))
	buffer8.WriteString(str1.String())
	buffer8.WriteString(str2.String())
	vv := buffer8.String()
	for n := 0; n < 9; n++ {
		buffer8.WriteString(vv)
	}

	return buffer8.String()
}

func Krandv2(size int, kind uint32) []byte {
	ikind, kinds, result := kind, [][]uint32{[]uint32{10, 48}, []uint32{26, 97}, []uint32{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	//rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			//ikind = rand.Intn(3)
			ikind = fastrand.Uint32n(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		//rand.Intn(scope)
		result[i] = uint8(base + fastrand.Uint32n(scope))
	}
	return result
}

type SS_fivecol struct {
	Ch_pronum int64
	Ch_bgjg   string
	Ch_khh    string
	Ch_tzm    string
	Ch_yswbh  string
	Ch_jyrq   string
}

func Rand_Line(chin chan *SS_fivecol) {

	//报告机构编码14bit
	//客户号7bit数字
	//大额交易特征码 4位
	//业务识别号 24bit,大写+数字
	//大额交易日期 年月日

	runNum := runtime.NumCPU()
	for n := 0; n < runNum; n++ {

		go func() {
			pronumtmp := time.Now().UnixNano() % 1000
			for {
				num := int(fastrand.Uint32n(10000))
				chin <- &SS_fivecol{
					Ch_pronum: pronumtmp,
					Ch_bgjg:   public.Map_jgbm[num],
					Ch_khh:    string(Krandv2(7, 0)),
					Ch_tzm:    string(Krandv2(4, 0)),
					Ch_yswbh:  string(Krandv2(24, 3)),
					Ch_jyrq:   Randate(),
				}
			}
		}()
	}

}

func Randate() string {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2020, 10, 0, 0, 0, 0, 0, time.UTC).Unix()
	detal := max - min
	sec := int64(fastrand.Uint32n(uint32(detal)) + uint32(min))
	return time.Unix(sec, 0).Format("20060102")
}
