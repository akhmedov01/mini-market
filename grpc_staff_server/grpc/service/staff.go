package service

import (
	"context"
	staff_service "staff/genproto"
	grpc_client "staff/grpc/client"
	"staff/packages/logger"
	"staff/storage"
)

type StaffService struct {
	logger  logger.LoggerI
	storage storage.StoregeI
	clients grpc_client.GrpcClientI
	staff_service.UnimplementedStaffServerServer
}

func NewStaffService(log logger.LoggerI, strg storage.StoregeI, grpcClients grpc_client.GrpcClientI) *StaffService {
	return &StaffService{
		logger:  log,
		storage: strg,
		clients: grpcClients,
	}
}

func (s *StaffService) Create(ctx context.Context, req *staff_service.CreateStaff) (*staff_service.IdRequest, error) {
	id, err := s.storage.Staff().CreateStaff(ctx, req)
	if err != nil {
		return nil, err
	}

	return &staff_service.IdRequest{Id: id}, nil
}

func (s *StaffService) Update(ctx context.Context, req *staff_service.Staff) (*staff_service.ResponseString, error) {
	str, err := s.storage.Staff().UpdateStaff(ctx, req)
	if err != nil {
		return nil, err
	}

	return &staff_service.ResponseString{Text: str}, nil
}

func (s *StaffService) Get(ctx context.Context, req *staff_service.IdRequest) (*staff_service.Staff, error) {
	staff, err := s.storage.Staff().GetStaff(ctx, req)
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *StaffService) GetAll(ctx context.Context, req *staff_service.GetAllStaffRequest) (*staff_service.GetAllStaffResponse, error) {
	staff, err := s.storage.Staff().GetAllStaff(ctx, req)
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *StaffService) Delete(ctx context.Context, req *staff_service.IdRequest) (*staff_service.ResponseString, error) {
	text, err := s.storage.Staff().DeleteStaff(ctx, req)
	if err != nil {
		return nil, err
	}

	return &staff_service.ResponseString{Text: text}, nil
}

func (s *StaffService) UpdateBalance(ctx context.Context, req *staff_service.UpdateBalanceRequest) (*staff_service.ResponseString, error) {
	text, err := s.storage.Staff().UpdateBalance(ctx, req)
	if err != nil {
		return nil, err
	}

	return &staff_service.ResponseString{Text: text}, nil
}

func (s *StaffService) GetByUsername(ctx context.Context, req *staff_service.RequestByUsername) (*staff_service.Staff, error) {
	staff, err := s.storage.Staff().GetByUsername(ctx, req)
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *StaffService) ChangePassword(ctx context.Context, req *staff_service.RequestByPassword) (*staff_service.ResponseString, error) {
	text, err := s.storage.Staff().ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return &staff_service.ResponseString{Text: text}, nil
}
