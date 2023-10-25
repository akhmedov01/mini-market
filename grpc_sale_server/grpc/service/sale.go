package service

import (
	"context"
	sale_service "sale/genproto"
	grpc_client "sale/grpc/client"
	"sale/packages/logger"
	"sale/storage"
)

type SaleService struct {
	logger  logger.LoggerI
	storage storage.StoregeI
	clients grpc_client.GrpcClientI
	sale_service.UnimplementedSaleServerServer
}

func NewSaleService(log logger.LoggerI, strg storage.StoregeI, grpcClients grpc_client.GrpcClientI) *SaleService {
	return &SaleService{
		logger:  log,
		storage: strg,
		clients: grpcClients,
	}
}

func (t *SaleService) Create(ctx context.Context, req *sale_service.CreateSale) (*sale_service.IdRequest, error) {
	id, err := t.storage.Sales().CreateSale(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.IdRequest{Id: id}, nil
}

func (t *SaleService) Update(ctx context.Context, req *sale_service.Sale) (*sale_service.ResponseString, error) {
	str, err := t.storage.Sales().UpdateSale(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.ResponseString{Text: str}, nil
}

func (t *SaleService) Get(ctx context.Context, req *sale_service.IdRequest) (*sale_service.Sale, error) {
	sale, err := t.storage.Sales().GetSale(ctx, req)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (t *SaleService) GetAll(ctx context.Context, req *sale_service.GetAllSaleRequest) (*sale_service.GetAllSaleResponse, error) {
	sales, err := t.storage.Sales().GetAllSale(ctx, req)
	if err != nil {
		return nil, err
	}

	return sales, nil
}

func (t *SaleService) Delete(ctx context.Context, req *sale_service.IdRequest) (*sale_service.ResponseString, error) {
	text, err := t.storage.Sales().DeleteSale(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.ResponseString{Text: text}, nil
}
