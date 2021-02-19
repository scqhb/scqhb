package services

import (
	"context"
)

type ProdService struct {
}

func (ss *ProdService) GetProdStock(context.Context, *ProdRequest) (*ProdResponse, error) {

	return &ProdResponse{ProdStock: 2000}, nil
}
func (ss *ProdService)  GetProdStocks(ctx context.Context, size *QuerySize) (*ProdResponseList, error){
	res:= []*ProdResponse{
		&ProdResponse{ProdStock:11},
		&ProdResponse{ProdStock:22},
		&ProdResponse{ProdStock:33},
		&ProdResponse{ProdStock:44},
		&ProdResponse{ProdStock:55},
	}
	return &ProdResponseList{Prodres:res},nil
}
