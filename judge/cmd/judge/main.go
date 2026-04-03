package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Yogdunana/yogduoj/judge/internal/config"
	pb "github.com/Yogdunana/yogduoj/judge/proto"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting YogduOJ Judge Service...")
	log.Printf("  Pool Size: %d", cfg.PoolSize)
	log.Printf("  gRPC Addr: %s", cfg.GRPCAddr)

	// Create gRPC listener
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", cfg.GRPCAddr, err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register judge service (stub - to be implemented)
	// pb.RegisterJudgeServiceServer(grpcServer, &service.JudgeService{})

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		grpcServer.GracefulStop()
	}()

	log.Printf("YogduOJ Judge Service is listening on %s", cfg.GRPCAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

	log.Println("Judge service stopped.")
}
