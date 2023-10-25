package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"staff/config"
	"staff/grpc"
	grpc_client "staff/grpc/client"

	"staff/packages/logger"
	"staff/storage/memory"
)

func main() {
	cfg := config.Load()
	lg := logger.NewLogger(cfg.Environment, "debug")
	strg, err := memory.NewStorage(context.Background(), *cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	clients, err := grpc_client.New(*cfg)
	if err != nil {
		log.Fatalf("failed to connect to services: %v", err)
	}
	s := grpc.SetUpServer(lg, strg, clients)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
