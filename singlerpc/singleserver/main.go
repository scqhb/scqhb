package singleserver

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"scdata/until"
	"strconv"

	"scdata/etcd/rpcserver"
	pb "scdata/singlerpc"

	"time"
)

/*const (
	port = "9.9.9.80:8003"
)
*/
type S struct {
	ServerAddress string
}

func (s *S) SingleRpc(ctx context.Context, in *pb.ClientRequest) (*pb.ServerResponse, error) {

	//log.Printf("收到客户端信息: %v", in.Requestmess)

	return &pb.ServerResponse{Responsemess: strconv.FormatBool(until.Filter.TestString(in.Requestmess))}, nil
}

/*
func SingleRpcMain() {
 	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	fmt.Println("启动single rpc 完成")
	s := grpc.NewServer()
	pb.RegisterSingleRpcBloomServiceServer(s, &S{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}



}
*/

///////////////////////////////////
const (
	key string = "grpcserverkey"
)

func SingleRpcMain(ip_port string) {
	flag.Parse()

	conf := &rpcserver.RpcServiceConf{
		Key:           key,
		ServerAddress: ip_port,
		Endpoints:     []string{"9.9.9.64:2379"},
		DialTimeout:   time.Second * 5,
	}

	server, err := rpcserver.NewRpcServer(conf, func(server *grpc.Server) {
		pb.RegisterSingleRpcBloomServiceServer(server, &S{conf.ServerAddress})
		//rpcfile.RegisterDemoServiceServer(server, &rpcserverimpl.DemoServiceServerimp{ServerAddress:conf.ServerAddress})
	})
	if err != nil {
		panic(err)
	}
	if err = server.Run(); err != nil {
		panic(err)
	}
}
