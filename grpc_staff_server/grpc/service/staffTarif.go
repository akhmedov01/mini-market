package service

import (
	"context"
	staff_service "staff/genproto"
	grpc_client "staff/grpc/client"
	"staff/packages/logger"
	"staff/storage"
)

type TarifService struct {
	logger  logger.LoggerI
	storage storage.StoregeI
	clients grpc_client.GrpcClientI
	staff_service.UnimplementedTarifServerServer
}

func NewTarifService(log logger.LoggerI, strg storage.StoregeI, grpcClients grpc_client.GrpcClientI) *TarifService {
	return &TarifService{
		logger:  log,
		storage: strg,
		clients: grpcClients,
	}
}

func (t *TarifService) Create(ctx context.Context, req *staff_service.CreateTarif) (*staff_service.IdRequest, error) {
	id, err := t.storage.StaffTarif().CreateStaffTarif(ctx, req)
	if err != nil {
		return nil, err
	}

	return &staff_service.IdRequest{Id: id}, nil
}

func (t *TarifService) Update(ctx context.Context, req *staff_service.Tarif) (*staff_service.ResponseString, error) {
	str, err := t.storage.StaffTarif().UpdateStaffTarif(ctx, req)
	if err != nil {
		return nil, err
	}

	return &staff_service.ResponseString{Text: str}, nil
}

func (t *TarifService) Get(ctx context.Context, req *staff_service.IdRequest) (*staff_service.Tarif, error) {
	staff, err := t.storage.StaffTarif().GetStaffTarif(ctx, req)
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (t *TarifService) GetAll(ctx context.Context, req *staff_service.GetAllTarifRequest) (*staff_service.GetAllTarifResponse, error) {
	staff, err := t.storage.StaffTarif().GetAllStaffTarif(ctx, req)
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (t *TarifService) Delete(ctx context.Context, req *staff_service.IdRequest) (*staff_service.ResponseString, error) {
	text, err := t.storage.StaffTarif().DeleteStaffTarif(ctx, req)
	if err != nil {
		return nil, err
	}

	return &staff_service.ResponseString{Text: text}, nil
}
