package grpc

import (
	branch_service "branch/genproto"
	grpc_client "branch/grpc/client"
	"branch/grpc/service"
	"branch/packages/logger"
	"branch/storage"

	"google.golang.org/grpc"
)

func SetUpServer(log logger.LoggerI, strg storage.StoregeI, grpcClient grpc_client.GrpcClientI) *grpc.Server {
	s := grpc.NewServer()
	branch_service.RegisterBranchServiceServer(s, service.NewBranchService(log, strg, grpcClient))
	return s
}
