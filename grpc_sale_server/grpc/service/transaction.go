package service

import (
	"context"
	sale_service "sale/genproto"
	grpc_client "sale/grpc/client"
	"sale/packages/logger"
	"sale/storage"
)

type TransactionService struct {
	logger  logger.LoggerI
	storage storage.StoregeI
	clients grpc_client.GrpcClientI
	sale_service.UnimplementedTransactionServerServer
}

func NewTransactionService(log logger.LoggerI, strg storage.StoregeI, grpcClients grpc_client.GrpcClientI) *TransactionService {
	return &TransactionService{
		logger:  log,
		storage: strg,
		clients: grpcClients,
	}
}

func (t *TransactionService) Create(ctx context.Context, req *sale_service.CreateTransaction) (*sale_service.IdRequest, error) {
	id, err := t.storage.Transaction().CreateStaffTransaction(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.IdRequest{Id: id}, nil
}

func (t *TransactionService) Update(ctx context.Context, req *sale_service.Transaction) (*sale_service.ResponseString, error) {
	str, err := t.storage.Transaction().UpdateStaffTransaction(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.ResponseString{Text: str}, nil
}

func (t *TransactionService) Get(ctx context.Context, req *sale_service.IdRequest) (*sale_service.Transaction, error) {
	transaction, err := t.storage.Transaction().GetStaffTransaction(ctx, req)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionService) GetAll(ctx context.Context, req *sale_service.GetAllTransactionRequest) (*sale_service.GetAllTransactionResponse, error) {
	transactions, err := t.storage.Transaction().GetAllStaffTransaction(ctx, req)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t *TransactionService) Delete(ctx context.Context, req *sale_service.IdRequest) (*sale_service.ResponseString, error) {
	text, err := t.storage.Transaction().DeleteStaffTransaction(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.ResponseString{Text: text}, nil
}
