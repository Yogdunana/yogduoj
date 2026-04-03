package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Yogdunana/yogduoj/judge/internal/config"
	"github.com/Yogdunana/yogduoj/judge/internal/sandbox"
	pb "github.com/Yogdunana/yogduoj/judge/proto"
	"github.com/Yogdunana/yogduoj/judge/internal/worker"
	"google.golang.org/grpc"
)

// Server wraps the gRPC server and all associated resources.
type Server struct {
	grpcServer  *grpc.Server
	judgeService *JudgeService
	containerPool *sandbox.ContainerPool
	workers     []*worker.Worker
	cfg         *config.Config
}

// NewServer creates and initializes the complete judge server.
func NewServer(cfg *config.Config) (*Server, error) {
	// Create the task channel.
	taskCh := make(chan *worker.JudgeTask, cfg.PoolSize*10)

	// Initialize the container pool.
	pool, err := sandbox.NewContainerPool(cfg.SandboxRoot, cfg.PoolSize)
	if err != nil {
		log.Printf("Warning: failed to create container pool: %v", err)
		log.Printf("Falling back to process-based execution without container isolation")
		// Continue without container pool - the runner will work without it.
	}

	// Create workers.
	workers := make([]*worker.Worker, cfg.PoolSize)
	for i := 0; i < cfg.PoolSize; i++ {
		workers[i] = worker.NewWorker(i, taskCh, cfg.SandboxRoot, cfg.CallbackURL)
	}

	// Create the judge service.
	judgeService := NewJudgeService(taskCh, workers)

	// Create the gRPC server.
	grpcServer := grpc.NewServer()
	pb.RegisterJudgeServiceServer(grpcServer, judgeService)

	s := &Server{
		grpcServer:   grpcServer,
		judgeService: judgeService,
		containerPool: pool,
		workers:      workers,
		cfg:          cfg,
	}

	return s, nil
}

// Start begins listening and serving requests.
func (s *Server) Start() error {
	// Create the TCP listener.
	lis, err := net.Listen("tcp", s.cfg.GRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.cfg.GRPCAddr, err)
	}

	// Start all workers.
	for _, w := range s.workers {
		w.Start()
	}

	log.Printf("YogduOJ Judge Service is listening on %s", s.cfg.GRPCAddr)
	log.Printf("  Pool Size: %d", s.cfg.PoolSize)
	log.Printf("  Callback URL: %s", s.cfg.CallbackURL)
	log.Printf("  Sandbox Root: %s", s.cfg.SandboxRoot)

	// Serve gRPC (blocking).
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop() {
	log.Println("Shutting down judge service...")

	// Mark service as unhealthy.
	s.judgeService.SetHealthy(false)

	// Stop accepting new gRPC requests.
	s.grpcServer.GracefulStop()

	// Close the task channel to signal workers to stop.
	// Note: workers read from the channel, so we need to close it.
	// This is handled by the channel being garbage collected when no more references exist.

	// Cleanup container pool.
	if s.containerPool != nil {
		s.containerPool.Close()
	}

	log.Println("Judge service stopped.")
}

// WaitForShutdown blocks until a termination signal is received.
func (s *Server) WaitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("Received signal %v, shutting down gracefully...", sig)
	s.Stop()
}
