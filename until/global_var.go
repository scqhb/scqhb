package until

import "sync"

var Ch1_ALL chan string = make(chan string, 100)
var Ch1_UPPER chan string = make(chan string, 100)
var Ch1_NUM chan string = make(chan string, 100)
var Ch1_NUM7 chan string = make(chan string, 100)
var Ch1_NUM4 chan string = make(chan string, 100)
var Ch2_PHTR2 chan *PHTR2 = make(chan *PHTR2, 100)
var Ch2_Xmlfile chan string = make(chan string, 100)
var D chan struct{} = make(chan struct{})
var Wg sync.WaitGroup
