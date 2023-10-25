// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: sale.proto

package sale_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SaleServerClient is the client API for SaleServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SaleServerClient interface {
	Create(ctx context.Context, in *CreateSale, opts ...grpc.CallOption) (*IdRequest, error)
	Update(ctx context.Context, in *Sale, opts ...grpc.CallOption) (*ResponseString, error)
	Get(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Sale, error)
	Delete(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*ResponseString, error)
	GetAll(ctx context.Context, in *GetAllSaleRequest, opts ...grpc.CallOption) (*GetAllSaleResponse, error)
}

type saleServerClient struct {
	cc grpc.ClientConnInterface
}

func NewSaleServerClient(cc grpc.ClientConnInterface) SaleServerClient {
	return &saleServerClient{cc}
}

func (c *saleServerClient) Create(ctx context.Context, in *CreateSale, opts ...grpc.CallOption) (*IdRequest, error) {
	out := new(IdRequest)
	err := c.cc.Invoke(ctx, "/sale_service.SaleServer/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *saleServerClient) Update(ctx context.Context, in *Sale, opts ...grpc.CallOption) (*ResponseString, error) {
	out := new(ResponseString)
	err := c.cc.Invoke(ctx, "/sale_service.SaleServer/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *saleServerClient) Get(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Sale, error) {
	out := new(Sale)
	err := c.cc.Invoke(ctx, "/sale_service.SaleServer/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *saleServerClient) Delete(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*ResponseString, error) {
	out := new(ResponseString)
	err := c.cc.Invoke(ctx, "/sale_service.SaleServer/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *saleServerClient) GetAll(ctx context.Context, in *GetAllSaleRequest, opts ...grpc.CallOption) (*GetAllSaleResponse, error) {
	out := new(GetAllSaleResponse)
	err := c.cc.Invoke(ctx, "/sale_service.SaleServer/GetAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SaleServerServer is the server API for SaleServer service.
// All implementations must embed UnimplementedSaleServerServer
// for forward compatibility
type SaleServerServer interface {
	Create(context.Context, *CreateSale) (*IdRequest, error)
	Update(context.Context, *Sale) (*ResponseString, error)
	Get(context.Context, *IdRequest) (*Sale, error)
	Delete(context.Context, *IdRequest) (*ResponseString, error)
	GetAll(context.Context, *GetAllSaleRequest) (*GetAllSaleResponse, error)
	mustEmbedUnimplementedSaleServerServer()
}

// UnimplementedSaleServerServer must be embedded to have forward compatible implementations.
type UnimplementedSaleServerServer struct {
}

func (UnimplementedSaleServerServer) Create(context.Context, *CreateSale) (*IdRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedSaleServerServer) Update(context.Context, *Sale) (*ResponseString, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedSaleServerServer) Get(context.Context, *IdRequest) (*Sale, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedSaleServerServer) Delete(context.Context, *IdRequest) (*ResponseString, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedSaleServerServer) GetAll(context.Context, *GetAllSaleRequest) (*GetAllSaleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedSaleServerServer) mustEmbedUnimplementedSaleServerServer() {}

// UnsafeSaleServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SaleServerServer will
// result in compilation errors.
type UnsafeSaleServerServer interface {
	mustEmbedUnimplementedSaleServerServer()
}

func RegisterSaleServerServer(s grpc.ServiceRegistrar, srv SaleServerServer) {
	s.RegisterService(&SaleServer_ServiceDesc, srv)
}

func _SaleServer_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSale)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SaleServerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sale_service.SaleServer/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SaleServerServer).Create(ctx, req.(*CreateSale))
	}
	return interceptor(ctx, in, info, handler)
}

func _SaleServer_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Sale)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SaleServerServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sale_service.SaleServer/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SaleServerServer).Update(ctx, req.(*Sale))
	}
	return interceptor(ctx, in, info, handler)
}

func _SaleServer_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SaleServerServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sale_service.SaleServer/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SaleServerServer).Get(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SaleServer_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SaleServerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sale_service.SaleServer/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SaleServerServer).Delete(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SaleServer_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllSaleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SaleServerServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sale_service.SaleServer/GetAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SaleServerServer).GetAll(ctx, req.(*GetAllSaleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SaleServer_ServiceDesc is the grpc.ServiceDesc for SaleServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SaleServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sale_service.SaleServer",
	HandlerType: (*SaleServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _SaleServer_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _SaleServer_Update_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _SaleServer_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _SaleServer_Delete_Handler,
		},
		{
			MethodName: "GetAll",
			Handler:    _SaleServer_GetAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sale.proto",
}