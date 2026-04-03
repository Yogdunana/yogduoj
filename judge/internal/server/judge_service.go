package server

import (
	"context"
	"log"
	"sync"
	"sync/atomic"

	pb "github.com/Yogdunana/yogduoj/judge/proto"
	"github.com/Yogdunana/yogduoj/judge/internal/worker"
)

// JudgeService implements the gRPC JudgeServiceServer interface.
type JudgeService struct {
	pb.UnimplementedJudgeServiceServer

	taskCh      chan *worker.JudgeTask
	workers     []*worker.Worker
	poolSize    int32
	runningCnt  int32
	queueLen    int32
	mu          sync.Mutex
	healthy     int32 // atomic: 1 = healthy, 0 = unhealthy
}

// NewJudgeService creates a new JudgeService with the given configuration.
func NewJudgeService(taskCh chan *worker.JudgeTask, workers []*worker.Worker) *JudgeService {
	s := &JudgeService{
		taskCh:   taskCh,
		workers:  workers,
		poolSize: int32(len(workers)),
		healthy:  1,
	}
	return s
}

// Submit receives a judge request, enqueues it for processing, and returns
// immediately with the submission ID. The actual result is sent asynchronously
// via the callback channel and HTTP callback.
func (s *JudgeService) Submit(ctx context.Context, req *pb.JudgeRequest) (*pb.JudgeResponse, error) {
	log.Printf("Received submit request: submission_id=%s, language=%s, test_cases=%d",
		req.SubmissionId, req.Language, len(req.TestCases))

	// Create a task with a done channel for synchronous response.
	done := make(chan *pb.JudgeResponse, 1)
	task := &worker.JudgeTask{
		Request: req,
		Done:    done,
	}

	// Enqueue the task.
	select {
	case s.taskCh <- task:
		atomic.AddInt32(&s.queueLen, 1)
	default:
		// Queue is full.
		log.Printf("Judge queue is full, rejecting submission %s", req.SubmissionId)
		return &pb.JudgeResponse{
			SubmissionId: req.SubmissionId,
			Verdict:      "SE",
			ErrorMessage: "Judge queue is full, please try again later",
		}, nil
	}

	// Wait for the result (with context timeout).
	select {
	case result := <-done:
		atomic.AddInt32(&s.queueLen, -1)
		return result, nil
	case <-ctx.Done():
		atomic.AddInt32(&s.queueLen, -1)
		log.Printf("Context cancelled while waiting for submission %s", req.SubmissionId)
		return &pb.JudgeResponse{
			SubmissionId: req.SubmissionId,
			Verdict:      "SE",
			ErrorMessage: "Request timed out",
		}, nil
	}
}

// GetStatus returns the current status of the judge service.
func (s *JudgeService) GetStatus(ctx context.Context, _ *pb.Empty) (*pb.JudgeStatusResponse, error) {
	// Count running workers.
	running := int32(0)
	for _, w := range s.workers {
		if w.IsRunning() {
			running++
		}
	}

	return &pb.JudgeStatusResponse{
		QueueSize:    atomic.LoadInt32(&s.queueLen),
		RunningCount: running,
		PoolSize:     s.poolSize,
		Healthy:      atomic.LoadInt32(&s.healthy) == 1,
	}, nil
}

// SetHealthy sets the health status of the service.
func (s *JudgeService) SetHealthy(healthy bool) {
	if healthy {
		atomic.StoreInt32(&s.healthy, 1)
	} else {
		atomic.StoreInt32(&s.healthy, 0)
	}
}
