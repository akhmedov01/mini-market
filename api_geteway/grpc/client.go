package grpc_client

import (
	"fmt"
	"main/config"
	branch_service "main/genproto/branch-server"
	staff_service "main/genproto/staff-server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClientI ...
type GrpcClientI interface {
	BranchService() branch_service.BranchServiceClient
	TarifService() staff_service.TarifServerClient
	StaffService() staff_service.StaffServerClient
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
			"branch_service": branch_service.NewBranchServiceClient(connSale),
			"staff_service":  staff_service.NewStaffServerClient(connSale),
			"tarif_service":  staff_service.NewTarifServerClient(connSale),
		},
	}, nil
}

func (g *GrpcClient) BranchService() branch_service.BranchServiceClient {
	return g.connections["branch_service"].(branch_service.BranchServiceClient)
}

func (g *GrpcClient) TarifService() staff_service.TarifServerClient {
	return g.connections["tarif_service"].(staff_service.TarifServerClient)
}

func (g *GrpcClient) StaffService() staff_service.StaffServerClient {
	return g.connections["staff_service"].(staff_service.StaffServerClient)
}
