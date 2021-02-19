package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"scdata/2021-02-17/services"
)

func main() {
	clientConn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer clientConn.Close()
	Prodclient := services.NewProdServiceClient(clientConn)
	ProdRes, err := Prodclient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 100})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ProdRes.ProdStock)
	fmt.Println("############")
	respons, err := Prodclient.GetProdStocks(context.Background(), &services.QuerySize{Size: 100})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(respons.Prodres)

}
