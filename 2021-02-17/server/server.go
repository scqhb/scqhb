package main

import (
	"google.golang.org/grpc"
	"scdata/2021-02-17/services"
)



func main(){
	rpcServer := grpc.NewServer()

services.RegisterProdServiceServer(rpcServer,new(services.ProdService))

}

