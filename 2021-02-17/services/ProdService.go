package services

import (
	"context"


)

type ProdService struct {


}

func(ss *ProdService)	GetProdStock(context.Context, *ProdRequest) (*ProdResponse, error){
	return &ProdResponse{ProdStock:2000},nil
}