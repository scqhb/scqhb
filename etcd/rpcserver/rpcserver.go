package rpcserver

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	grpc "google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	regis8 "scdata/etcd/register"
	"syscall"
	"time"
)

type (
	RpcServiceFunc func(server *grpc.Server)

	RpcService struct {
		register       *regis8.Register
		RpcServiceFunc RpcServiceFunc
	}

	RpcServiceConf struct {
		Key           string
		ServerAddress string   //grpc 服务地址
		Endpoints     []string //etcd地址,9.9.9.64:2379
		DialTimeout   time.Duration
	}
)

func NewRpcServer(conf *RpcServiceConf, rpcServiceFunc RpcServiceFunc) (*RpcService, error) {
	//创建etcd链接
	var client3 *clientv3.Client
	var err error
	if client3, err = clientv3.New(clientv3.Config{Endpoints: conf.Endpoints, DialTimeout: conf.DialTimeout}); err != nil {
		fmt.Println("connect etcd err:", err)
		panic(err)
		return nil, err
	}

	register := regis8.NewRegister(conf.Key, client3, conf.ServerAddress) //返回一个结构体
	/*Register struct {
		Key           string
		cli3          *clientv3.Client
		serverAddress string
		stop          chan bool
		interval      time.Duration
		leaseTime     int64
	}*/
	return &RpcService{
		register:       register,
		RpcServiceFunc: rpcServiceFunc,
	}, nil
}

func (s *RpcService) Run(serverOptions ...grpc.ServerOption) error {

	listen, err := net.Listen("tcp", s.register.GetServerAddress()) //监听rpc服务
	if err != nil {
		fmt.Println("rpc server listen failed!!!", s.register.GetServerAddress())
		return err
	}
	log.Printf("Rpc server listen at: %s", s.register.GetServerAddress())
	s.register.Reg()
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxx")

	server := grpc.NewServer(serverOptions...)
	s.RpcServiceFunc(server)

	s.deadNotify()

	if err := server.Serve(listen); err != nil {
		return err
	}

	return nil
}

func (s *RpcService) deadNotify() error {
	fmt.Println("deadmotify begin....")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		log.Printf("signal.notify   system os notefiy:%v", <-ch)
		s.register.UnReg()
		os.Exit(1)
	}()
	fmt.Println("deadmotify end....")

	return nil
}
