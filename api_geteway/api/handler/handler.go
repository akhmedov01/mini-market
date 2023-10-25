package handler

import (
	"main/config"
	grpc_client "main/grpc"
	"main/packages/logger"
	"main/storage"
)

type Handler struct {
	cfg        config.Config
	log        logger.LoggerI
	red        storage.CacheI
	grpcClient grpc_client.GrpcClientI
}

func NewHandler(conf config.Config, loger logger.LoggerI, redis storage.CacheI, client grpc_client.GrpcClientI) *Handler {
	return &Handler{cfg: conf, log: loger, red: redis, grpcClient: client}
}
