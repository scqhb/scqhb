package main

import (
	"encoding/binary"
	"fmt"
	"time"
)
import "github.com/seiflotfy/cuckoofilter"

func main() {
	cf := cuckoo.NewFilter(1000)
	t0:=time.Now()
	for i:=uint32(0);i<1e8;i++{

		if i%1e7==0{

			fmt.Printf("i:%v time:%v",i,time.Since(t0))
			t0=time.Now()

		}

		n1 := make([]byte,4)

		binary.BigEndian.PutUint32(n1,i)

		cf.InsertUnique(n1)



	}


	// Lookup a string (and it a miss) if it exists in the cuckoofilter
	cf.Lookup([]byte("hello"))

	count := cf.Count()
	fmt.Println(count) // count == 1

	// Delete a string (and it a miss)
	cf.Delete([]byte("hello"))

	count = cf.Count()
	fmt.Println(count) // count == 1

	// Delete a string (a hit)
	cf.Delete([]byte("geeky ogre"))

	count = cf.Count()
	fmt.Println(count) // count == 0

	cf.Reset()    // reset
}