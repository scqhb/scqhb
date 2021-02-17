package main

import (
	"errors"
	"fmt"
	"os"
)

//
type Queue struct {
	maxSize int
	array [5]int
	front int
	rear int
}
func(this *Queue) AddQueue(val int)(err error){
	if this.rear==this.maxSize-1{
		return errors.New("queur full")
	}
	this.rear++
	this.array[this.rear]=val
	return
}
func(this *Queue)ShowQueue(){
	for i:=this.front+1;i<=this.rear;i++{
		fmt.Printf("arrar[%d]=%d\t",i,this.array[i])
	}
}

func(this *Queue) GetQueue()(val int,err error){
	if this.rear==this.front{
		return -1,errors.New("query empty")
	}
	this.front++
	val=this.array[this.front]
	return val,nil


}


func main(){
	queue:=&Queue{
		maxSize: 5,
		array:   [5]int{},
		front:   -1,
		rear:    -1,
	}
	var key string
	var val int
	for {
		fmt.Println("1 putllll lslssllssvalues")
		fmt.Println("2.get data")
		fmt.Println("3.show queue")
		fmt.Println("4.exit queue")
		fmt.Scanln(&key)
		switch key {
		case "add":
			fmt.Println(":")
			fmt.Scanln(&val)

			err := queue.AddQueue(val)
			if err != nil {
				fmt.Println(err.Error())

			} else {
				fmt.Println("sucess")

			}
		case "get":
			fmt.Println("get")
			getQueue, err := queue.GetQueue()
			if err!=nil{
				fmt.Println(err.Error())
			}else {
				fmt.Println("from queue value:",getQueue)
			}

		case "show":
			fmt.Println("show ....")
			queue.ShowQueue()
		case "exit":
			os.Exit(0)
		}

	}




}
