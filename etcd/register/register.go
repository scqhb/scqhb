package regis8

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/clientv3"
	"log"

	"time"
)

//定义数据类型
type (
	Register struct {
		Key           string
		cli3          *clientv3.Client
		serverAddress string
		stop          chan bool
		interval      time.Duration
		leaseTime     int64
	}
)

//构造体函数
func NewRegister(key string, cli3 *clientv3.Client, serverAddress string) *Register {
	return &Register{
		Key:           key, //etcd key值
		cli3:          cli3,
		serverAddress: serverAddress, //RPC服务地址
		stop:          make(chan bool, 1),
		interval:      time.Second * 1,
		leaseTime:     20, //租约过期时间
	}
}

//将key put进入etcd
func (r *Register) Reg() {
	keyname := r.makeKeyName() //生成key和serveraddress字符串
	go func() {
		t := time.NewTicker(r.interval) //定时往通道发送信息
		for {
			//创建一个新的租约
			lgs, err := r.cli3.Grant(context.TODO(), r.leaseTime) //生成租约
			if err != nil {
				panic(err)
			}
			if _, err := r.cli3.Get(context.TODO(), keyname); err != nil { //在etcd中获取key

				if err == rpctypes.ErrKeyNotFound { //如果不存在key,就插入key,key带有租约
					if _, err := r.cli3.Put(
						context.TODO(),
						keyname, r.serverAddress,
						clientv3.WithLease(lgs.ID)); err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			} else {
				if _, err := r.cli3.Put(
					context.TODO(), keyname,
					r.serverAddress,
					clientv3.WithLease(lgs.ID)); err != nil {
					panic(err)
				}
			}
			select {
			case tt := <-t.C:
				log.Println("ticker time:", tt.Format("2006-01-02 15:04:05"))
			case <-r.stop:
				return
			}
		}
	}()
}

//获取服务地址
func (r *Register) GetServerAddress() string {
	return r.serverAddress
}

//生成key
func (r *Register) makeKeyName() string {
	keyname := fmt.Sprintf("%s-%s", r.Key, r.serverAddress)
	fmt.Println("makekey:", keyname)
	return keyname
}

//unreg取消注册
func (r *Register) UnReg() {
	r.stop <- true
	keyname := r.makeKeyName()
	r.stop = make(chan bool, 1) //为了防止多线程下面死锁
	if _, err := r.cli3.Delete(context.TODO(), keyname); err != nil {
		panic(err)
	} else {
		log.Printf("%s unres sucess", keyname)
	}
	return

}
