package grpc_client

import (
	"branch/config"
	sale_service "branch/genproto"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClientI ...
type GrpcClientI interface {
	BranchService() sale_service.BranchServiceClient
}

// GrpcClient ...
type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

// New ...
func New(cfg config.Config) (*GrpcClient, error) {
	connSale, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.SaleServiceHost, cfg.SaleServisePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("courier service dial host: %s port:%s err: %s",
			cfg.SaleServiceHost, cfg.SaleServisePort, err)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"sale_service": sale_service.NewBranchServiceClient(connSale),
		},
	}, nil
}

func (g *GrpcClient) BranchService() sale_service.BranchServiceClient {
	return g.connections["sale_service"].(sale_service.BranchServiceClient)
}
