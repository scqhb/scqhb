package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"scdata/2021-02-17/services"
)

func main() {
	rpcServer := grpc.NewServer()

	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	rpcServer.Serve(listen)

}
