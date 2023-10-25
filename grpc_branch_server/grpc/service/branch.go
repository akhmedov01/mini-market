package service

import (
	branch_service "branch/genproto"
	grpc_client "branch/grpc/client"
	"branch/packages/logger"
	"branch/storage"
	"context"
)

type BranchService struct {
	logger  logger.LoggerI
	storage storage.StoregeI
	clients grpc_client.GrpcClientI
	branch_service.UnimplementedBranchServiceServer
}

func NewBranchService(log logger.LoggerI, strg storage.StoregeI, grpcClients grpc_client.GrpcClientI) *BranchService {
	return &BranchService{
		logger:  log,
		storage: strg,
		clients: grpcClients,
	}
}

func (b *BranchService) Create(ctx context.Context, req *branch_service.CreateBranch) (*branch_service.IdReqRes, error) {
	id, err := b.storage.Branch().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &branch_service.IdReqRes{Id: id}, nil
}

func (b *BranchService) Update(ctx context.Context, req *branch_service.Branch) (*branch_service.ResponseString, error) {
	str, err := b.storage.Branch().Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return &branch_service.ResponseString{Text: str}, nil
}

func (b *BranchService) Get(ctx context.Context, req *branch_service.IdReqRes) (*branch_service.Branch, error) {
	branch, err := b.storage.Branch().GetBranch(ctx, req)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

func (b *BranchService) GetAll(ctx context.Context, req *branch_service.GetAllBranchRequest) (*branch_service.GetAllBranchResponse, error) {
	branch, err := b.storage.Branch().GetAllBranch(ctx, req)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

func (b *BranchService) Delete(ctx context.Context, req *branch_service.IdReqRes) (*branch_service.ResponseString, error) {
	text, err := b.storage.Branch().DeleteBranch(ctx, req)
	if err != nil {
		return nil, err
	}

	return &branch_service.ResponseString{Text: text}, nil
}
