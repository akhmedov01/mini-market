package grpc

import (
	staff_service "staff/genproto"
	grpc_client "staff/grpc/client"
	"staff/grpc/service"
	"staff/packages/logger"
	"staff/storage"

	"google.golang.org/grpc"
)

func SetUpServer(log logger.LoggerI, strg storage.StoregeI, grpcClient grpc_client.GrpcClientI) *grpc.Server {
	s := grpc.NewServer()
	staff_service.RegisterTarifServerServer(s, service.NewTarifService(log, strg, grpcClient))
	return s
}
