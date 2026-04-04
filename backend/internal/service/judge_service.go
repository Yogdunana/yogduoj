package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
	pb "github.com/Yogdunana/yogduoj/judge/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// JudgeNotifyFunc is a callback function type used to notify about judge status updates.
// This breaks the import cycle between service and handler packages.
type JudgeNotifyFunc func(submissionID uint, data interface{})

// JudgeService handles code judging operations.
type JudgeService interface {
	Submit(ctx context.Context, submission *model.Submission) error
	GetJudgeStatus(ctx context.Context, submissionID uint) (*model.Submission, error)
}

type judgeService struct {
	judgeClient   pb.JudgeServiceClient
	grpcConn      *grpc.ClientConn
	problemRepo   repository.ProblemRepository
	submissionRepo repository.SubmissionRepository
	logger        *zap.Logger
	timeout       time.Duration
	notifyFunc    JudgeNotifyFunc
}

// NewJudgeService creates a new JudgeService that connects to the Judge gRPC server.
func NewJudgeService(grpcAddr string, timeoutSeconds int, problemRepo repository.ProblemRepository, submissionRepo repository.SubmissionRepository) JudgeService {
	if grpcAddr == "" {
		grpcAddr = "judge:50051"
	}
	if timeoutSeconds <= 0 {
		timeoutSeconds = 30
	}

	timeout := time.Duration(timeoutSeconds) * time.Second

	// Connect to the judge gRPC server with insecure credentials (internal network).
	// In production, consider using TLS if the judge is on a different network.
	conn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		// Log the error but don't panic - the service will retry on each submission
		fmt.Fprintf(os.Stderr, "[judge] WARNING: failed to connect to judge gRPC server at %s: %v\n", grpcAddr, err)
		conn = nil
	}

	var client pb.JudgeServiceClient
	if conn != nil {
		client = pb.NewJudgeServiceClient(conn)
	}

	return &judgeService{
		judgeClient:    client,
		grpcConn:       conn,
		problemRepo:    problemRepo,
		submissionRepo: submissionRepo,
		logger:         zap.NewNop(),
		timeout:        timeout,
	}
}

// SetLogger allows injecting a logger after construction.
func (s *judgeService) SetLogger(logger *zap.Logger) {
	s.logger = logger
}

// SetNotifyFunc sets the callback function for judge status notifications.
// This should be called during application initialization to wire up WebSocket notifications.
func (s *judgeService) SetNotifyFunc(fn JudgeNotifyFunc) {
	s.notifyFunc = fn
}

// Close cleans up the gRPC connection.
func (s *judgeService) Close() {
	if s.grpcConn != nil {
		_ = s.grpcConn.Close()
	}
}

// Submit sends a submission to the judge gRPC server for evaluation.
// It loads the problem's test data, reads the source code, builds a JudgeRequest,
// calls the gRPC service, and updates the submission with the results.
// The judge call runs asynchronously in a goroutine so it doesn't block the HTTP response.
func (s *judgeService) Submit(ctx context.Context, submission *model.Submission) error {
	// Run the judging in a goroutine so the HTTP response returns immediately.
	// The client will get updates via WebSocket.
	go s.processSubmission(submission)

	return nil
}

// notify sends a judge status update through the registered callback.
func (s *judgeService) notify(submissionID uint, data interface{}) {
	if s.notifyFunc != nil {
		s.notifyFunc(submissionID, data)
	}
}

// processSubmission does the actual work of calling the judge service.
func (s *judgeService) processSubmission(submission *model.Submission) {
	bgCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	// Update status to "judging"
	_ = s.submissionRepo.UpdateSubmissionStatus(bgCtx, submission.ID, "judging")

	// Notify WebSocket clients about status change
	s.notify(submission.ID, map[string]interface{}{
		"submission_id": submission.ID,
		"judge_result":  "judging",
	})

	// Load the problem to get time/memory limits and test data
	problem, err := s.problemRepo.GetByID(bgCtx, submission.ProblemID)
	if err != nil {
		s.handleJudgeError(submission, fmt.Sprintf("failed to load problem: %v", err))
		return
	}

	// Load test data for the problem
	testDataList, err := s.problemRepo.GetTestData(bgCtx, submission.ProblemID)
	if err != nil {
		s.handleJudgeError(submission, fmt.Sprintf("failed to load test data: %v", err))
		return
	}

	// Read the source code from the file
	code, err := os.ReadFile(submission.CodePath)
	if err != nil {
		s.handleJudgeError(submission, fmt.Sprintf("failed to read source code: %v", err))
		return
	}

	// Build the gRPC JudgeRequest
	req := &pb.JudgeRequest{
		SubmissionId:     strconv.FormatUint(uint64(submission.ID), 10),
		Language:         submission.Language,
		Code:             string(code),
		TimeLimitMs:      int64(problem.TimeLimitMs),
		MemoryLimitKb:    int64(problem.MemoryLimitKb),
		SpecialJudgePath: problem.SpecialJudgePath,
		CtfFlag:          problem.CTFFlag,
	}

	// Convert test data to proto TestCase format
	for _, td := range testDataList {
		input, err := os.ReadFile(td.InputFile)
		if err != nil {
			s.logger.Warn("failed to read test input file",
				zap.Uint("test_data_id", td.ID),
				zap.String("file", td.InputFile),
				zap.Error(err),
			)
			continue
		}

		expectedOutput, err := os.ReadFile(td.OutputFile)
		if err != nil {
			s.logger.Warn("failed to read test output file",
				zap.Uint("test_data_id", td.ID),
				zap.String("file", td.OutputFile),
				zap.Error(err),
			)
			continue
		}

		req.TestCases = append(req.TestCases, &pb.TestCase{
			Input:          string(input),
			ExpectedOutput: string(expectedOutput),
			ScoreWeight:    int32(td.ScoreWeight),
		})
	}

	// Check if we have a gRPC client
	if s.judgeClient == nil {
		s.logger.Warn("judge gRPC client not available",
			zap.String("error", "judge service is not available"),
		)
		s.handleJudgeError(submission, "judge service is not available")
		return
	}

	// Call the judge gRPC service
	resp, err := s.judgeClient.Submit(bgCtx, req)
	if err != nil {
		s.handleJudgeError(submission, fmt.Sprintf("judge service call failed: %v", err))
		return
	}

	// Build judge detail from case results
	judgeDetail, err := buildJudgeDetail(resp)
	if err != nil {
		s.logger.Warn("failed to serialize judge detail",
			zap.Uint("submission_id", submission.ID),
			zap.Error(err),
		)
	}

	// Update the submission with judge results
	if err := s.submissionRepo.UpdateJudgeResult(
		bgCtx,
		submission.ID,
		resp.Verdict,
		float64(resp.TotalScore),
		int(resp.TimeUsedMs),
		int(resp.MemoryUsedKb),
		resp.ErrorMessage,
		judgeDetail,
	); err != nil {
		s.logger.Error("failed to update judge result in database",
			zap.Uint("submission_id", submission.ID),
			zap.Error(err),
		)
		return
	}

	// Notify WebSocket clients about the final result
	s.notify(submission.ID, map[string]interface{}{
		"submission_id":  submission.ID,
		"judge_result":   resp.Verdict,
		"judge_score":    resp.TotalScore,
		"time_used_ms":   resp.TimeUsedMs,
		"memory_used_kb": resp.MemoryUsedKb,
		"error_message":  resp.ErrorMessage,
	})

	s.logger.Info("submission judged successfully",
		zap.Uint("submission_id", submission.ID),
		zap.String("verdict", resp.Verdict),
		zap.Int32("score", resp.TotalScore),
		zap.Int64("time_ms", resp.TimeUsedMs),
		zap.Int64("memory_kb", resp.MemoryUsedKb),
	)
}

// handleJudgeError updates a submission with an error verdict and notifies WebSocket clients.
func (s *judgeService) handleJudgeError(submission *model.Submission, errMsg string) {
	bgCtx := context.Background()

	s.logger.Error("judge error",
		zap.Uint("submission_id", submission.ID),
		zap.String("error", errMsg),
	)

	_ = s.submissionRepo.UpdateJudgeResult(
		bgCtx,
		submission.ID,
		"SE", // System Error
		0,
		0,
		0,
		errMsg,
		"",
	)

	s.notify(submission.ID, map[string]interface{}{
		"submission_id": submission.ID,
		"judge_result":  "SE",
		"error_message": errMsg,
	})
}

// GetJudgeStatus returns the current status of a submission.
func (s *judgeService) GetJudgeStatus(ctx context.Context, submissionID uint) (*model.Submission, error) {
	return s.submissionRepo.GetByID(ctx, submissionID)
}

// buildJudgeDetail creates a JSON string with per-test-case results from the gRPC response.
func buildJudgeDetail(resp *pb.JudgeResponse) (string, error) {
	type caseResult struct {
		Index        int32  `json:"index"`
		Verdict      string `json:"verdict"`
		TimeUsedMs   int64  `json:"time_used_ms"`
		MemoryUsedKb int64  `json:"memory_used_kb"`
		ErrorMessage string `json:"error_message,omitempty"`
		Score        int32  `json:"score"`
	}

	detail := struct {
		Verdict      string       `json:"verdict"`
		TotalScore   int32        `json:"total_score"`
		TimeUsedMs   int64        `json:"time_used_ms"`
		MemoryUsedKb int64        `json:"memory_used_kb"`
		CaseResults  []caseResult `json:"case_results"`
	}{
		Verdict:      resp.Verdict,
		TotalScore:   resp.TotalScore,
		TimeUsedMs:   resp.TimeUsedMs,
		MemoryUsedKb: resp.MemoryUsedKb,
	}

	for _, cr := range resp.CaseResults {
		detail.CaseResults = append(detail.CaseResults, caseResult{
			Index:        cr.Index,
			Verdict:      cr.Verdict,
			TimeUsedMs:   cr.TimeUsedMs,
			MemoryUsedKb: cr.MemoryUsedKb,
			ErrorMessage: cr.ErrorMessage,
			Score:        cr.Score,
		})
	}

	data, err := json.Marshal(detail)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// judgeServiceOnce ensures the notify function is set only once.
var (
	judgeServiceOnce sync.Once
	notifyFuncStore  JudgeNotifyFunc
)

// SetGlobalNotifyFunc sets the global judge notification callback.
// This is called from the handler package during initialization to wire up WebSocket notifications.
func SetGlobalNotifyFunc(fn JudgeNotifyFunc) {
	judgeServiceOnce.Do(func() {
		notifyFuncStore = fn
	})
}

// GetGlobalNotifyFunc returns the global judge notification callback.
func GetGlobalNotifyFunc() JudgeNotifyFunc {
	return notifyFuncStore
}
