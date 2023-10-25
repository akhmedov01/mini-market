package service

import (
	"context"
	sale_service "sale/genproto"
	grpc_client "sale/grpc/client"
	"sale/packages/logger"
	"sale/storage"
)

type BranchTransactionService struct {
	logger  logger.LoggerI
	storage storage.StoregeI
	clients grpc_client.GrpcClientI
	sale_service.UnimplementedBranchTransactionServerServer
}

func NewBranchTransactionService(log logger.LoggerI, strg storage.StoregeI, grpcClients grpc_client.GrpcClientI) *BranchTransactionService {
	return &BranchTransactionService{
		logger:  log,
		storage: strg,
		clients: grpcClients,
	}
}

func (t *BranchTransactionService) Create(ctx context.Context, req *sale_service.CreateBranchTransaction) (*sale_service.IdRequest, error) {
	id, err := t.storage.BranchTransaction().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.IdRequest{Id: id}, nil
}

func (t *BranchTransactionService) Update(ctx context.Context, req *sale_service.BranchTransaction) (*sale_service.ResponseString, error) {
	str, err := t.storage.BranchTransaction().Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.ResponseString{Text: str}, nil
}

func (t *BranchTransactionService) Get(ctx context.Context, req *sale_service.IdRequest) (*sale_service.BranchTransaction, error) {
	branchTransaction, err := t.storage.BranchTransaction().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return branchTransaction, nil
}

func (t *BranchTransactionService) GetAll(ctx context.Context, req *sale_service.GetAllBranchTransactionRequest) (*sale_service.GetAllBranchTransactionResponse, error) {
	branchTransactions, err := t.storage.BranchTransaction().GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return branchTransactions, nil
}

func (t *BranchTransactionService) Delete(ctx context.Context, req *sale_service.IdRequest) (*sale_service.ResponseString, error) {
	text, err := t.storage.BranchTransaction().Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	return &sale_service.ResponseString{Text: text}, nil
}
