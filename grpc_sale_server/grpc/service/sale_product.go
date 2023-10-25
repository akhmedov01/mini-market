package service

import (
	"context"
	sale_service "sale/genproto"
	grpc_client "sale/grpc/client"
	"sale/packages/logger"
	"sale/storage"
)

type SaleProductService struct {
	logger  logger.LoggerI
	storage storage.StoregeI
	clients grpc_client.GrpcClientI
	sale_service.UnimplementedSaleProductServerServer
}

func NewSaleProductService(log logger.LoggerI, strg storage.StoregeI, grpcClients grpc_client.GrpcClientI) *SaleProductService {
	return &SaleProductService{
		logger:  log,
		storage: strg,
		clients: grpcClients,
	}
}

func (t *SaleProductService) Create(ctx context.Context, req *sale_service.CreateSaleProduct) (*sale_service.IdRequest, error) {
	id, err := t.storage.SaleProduct().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.IdRequest{Id: id}, nil
}

func (t *SaleProductService) Update(ctx context.Context, req *sale_service.SaleProduct) (*sale_service.ResponseString, error) {
	str, err := t.storage.SaleProduct().Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.ResponseString{Text: str}, nil
}

func (t *SaleProductService) Get(ctx context.Context, req *sale_service.IdRequest) (*sale_service.SaleProduct, error) {
	saleProduct, err := t.storage.SaleProduct().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return saleProduct, nil
}

func (t *SaleProductService) GetAll(ctx context.Context, req *sale_service.GetAllSaleProductRequest) (*sale_service.GetAllSaleProductResponse, error) {
	saleProducts, err := t.storage.SaleProduct().GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return saleProducts, nil
}

func (t *SaleProductService) Delete(ctx context.Context, req *sale_service.IdRequest) (*sale_service.ResponseString, error) {
	text, err := t.storage.SaleProduct().Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.ResponseString{Text: text}, nil
}
