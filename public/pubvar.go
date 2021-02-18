package public

import (
	bf "github.com/willf/bloom"
	"hash/crc32"
	"sync"
)

var Map_jgbm map[int]string = make(map[int]string)

const ArrySize = 5

var ArryBloomS [ArrySize]*bf.BloomFilter
var ArryChanData [ArrySize]chan []byte
var CrcTable = crc32.MakeTable(crc32.IEEE)
var ArryChan_select [ArrySize]chan []byte
var Exitchan chan interface{} = make(chan interface{}, ArrySize)
var FileNum int = 0
var Glock sync.Mutex

func GetFilenum() int {
	Glock.Lock()
	FileNum++
	tmp := FileNum
	Glock.Unlock()
	return tmp

}
