package grpc

import (
	"flag"
	"google.golang.org/grpc"
	"io"
	"log"
	"scdata/etcd/rpcserver"
	proto "scdata/grpc" // 根据proto文件自动生成的代码
	"scdata/until"
	"time"
	// 	"strconv"
)

// Streamer 服务端
type Streamer struct {
	ServerAddress string
}

// BidStream 实现了 ChatServer 接口中定义的 BidStream 方法
func (s *Streamer) BloomStream(stream proto.Chat_BloomStreamServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			log.Println("收到客户端通过context发出的终止信号")
			return ctx.Err()
		default:
			// 接收从客户端发来的消息
			recyin, err := stream.Recv()
			if err == io.EOF {
				log.Println("客户端发送的数据流结束")

				return nil
			}
			if err != nil {
				log.Println("接收数据出错:", err)
				return err
			}
			// 如果接收正常，则根据接收到的 字符串 执行相应的指令
			switch recyin.Input {
			case "stop":
				log.Println("收到'结束对话'指令")
				if err := stream.Send(&proto.Response{Output: "收到结束指令"}); err != nil {
					return err
				}
				// 收到结束指令时，通过 return nil 终止双向数据流
				return nil
				/*			case "sss"://流返回
							log.Println("收到'返回数据流'指令")
							// 收到 收到'返回数据流'指令， 连续返回 10 条数据
							for i := 0; i < 10; i++ {
								if err := stream.Send(&proto.Response{Output: "数据流 #" + strconv.Itoa(i)}); err != nil {
									return err
								}
							}*/
			default:
				// 缺省情况下， 返回 '服务端返回: ' + 输入信息
				if until.Filter.TestString(recyin.Input) {
					//	fmt.Printf("false positive data:%v\n", recyin.Input)
					if err := stream.Send(&proto.Response{Output: recyin.Input}); err != nil {
						log.Printf("send false positive data err:%v  data:%v\n", err, recyin.Input)
						return err
					}
				}
			}
		}
	}
}

const (
	key string = "grpcstreamserverkey"
)

func BloomServiceGrpcStream(ip_port string) {
	//ip_port="9.9.9.80:3003"
	log.Println("grpc stream 启动服务端...")
	/*	server := grpc.NewServer()
		// 注册 ChatServer
		proto.RegisterChatServer(server, &Streamer{})

		address, err := net.Listen("tcp", ":3000")
		if err != nil {
			panic(err)
		}
		if err := server.Serve(address); err != nil {
			panic(err)
		}
	*/
	////////////////###############
	flag.Parse()

	conf := &rpcserver.RpcServiceConf{
		Key:           key,
		ServerAddress: ip_port,
		Endpoints:     []string{"etcserver01:2379"},
		DialTimeout:   time.Second * 5,
	}

	server, err := rpcserver.NewRpcServer(conf, func(server *grpc.Server) {
		proto.RegisterChatServer(server, &Streamer{conf.ServerAddress})
	})
	if err != nil {
		panic(err)
	}
	if err = server.Run(); err != nil {
		panic(err)
	}
}

/*

const (
	key string = "grpcserverkey"
)



func SingleRpcMain(port *int) {
	flag.Parse()

	conf := &rpcserver.RpcServiceConf{
		Key:           key,
		ServerAddress: fmt.Sprintf("9.9.9.80:%d", *port),
		Endpoints:     []string{"9.9.9.64:2379"},
		DialTimeout:time.Second*5,
	}

	server, err := rpcserver.NewRpcServer(conf, func(server *grpc.Server) {
		pb.RegisterSingleRpcBloomServiceServer(server,&S{conf.ServerAddress})
		//rpcfile.RegisterDemoServiceServer(server, &rpcserverimpl.DemoServiceServerimp{ServerAddress:conf.ServerAddress})
	})
	if err != nil {
		panic(err)
	}
	if err = server.Run(); err != nil {
		panic(err)
	}
}

*/
