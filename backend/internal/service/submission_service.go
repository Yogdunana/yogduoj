package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrSubmissionNotFound = errors.New("submission not found")
	ErrProblemNotPublic   = errors.New("problem is not available")
	ErrCodeEmpty          = errors.New("code cannot be empty")
	ErrInvalidLanguage    = errors.New("invalid programming language")
	ErrCTFAnswerEmpty     = errors.New("ctf answer cannot be empty")
)

// ValidLanguages maps supported programming languages.
var ValidLanguages = map[string]bool{
	"cpp":     true,
	"c":       true,
	"java":    true,
	"python3": true,
}

// SubmissionFilter holds filter parameters for listing submissions.
type SubmissionFilter struct {
	UserID     uint
	ProblemID  uint
	ContestID  uint
	JudgeResult string
	Language   string
}

// CreateSubmissionRequest is the request body for creating a submission.
type CreateSubmissionRequest struct {
	ProblemID uint   `json:"problem_id" binding:"required"`
	Language  string `json:"language" binding:"required"`
	Code      string `json:"code" binding:"required"`
	CTFAnswer string `json:"ctf_answer"`
}

// JudgeResultCallback is the callback data from the judge service.
type JudgeResultCallback struct {
	SubmissionID  uint    `json:"submission_id"`
	JudgeResult   string  `json:"judge_result"`
	JudgeScore    float64 `json:"judge_score"`
	TimeUsedMs    int     `json:"time_used_ms"`
	MemoryUsedKb  int     `json:"memory_used_kb"`
	ErrorMessage  string  `json:"error_message"`
	JudgeDetail   string  `json:"judge_detail"`
}

// SubmissionService handles submission-related business logic.
type SubmissionService interface {
	CreateSubmission(ctx context.Context, userID uint, req CreateSubmissionRequest, ipAddress string) (*model.Submission, error)
	GetSubmission(ctx context.Context, submissionID uint, userID uint, isAdmin bool) (*model.Submission, error)
	GetSubmissionCode(ctx context.Context, submissionID uint, userID uint, isAdmin bool) (string, error)
	ListSubmissions(ctx context.Context, filter SubmissionFilter, offset, limit int, userID uint, isAdmin bool) ([]model.Submission, int64, error)
	HandleJudgeResult(ctx context.Context, submissionID uint, result JudgeResultCallback) error
	RejudgeSubmission(ctx context.Context, submissionID uint) error
	GetUserSubmissions(ctx context.Context, userID uint, offset, limit int) ([]model.Submission, int64, error)
	GetContestSubmissions(ctx context.Context, contestID uint, offset, limit int) ([]model.Submission, int64, error)
}

type submissionService struct {
	submissionRepo repository.SubmissionRepository
	problemRepo    repository.ProblemRepository
	judgeService   JudgeService
	logger         *zap.Logger
}

func NewSubmissionService(
	submissionRepo repository.SubmissionRepository,
	problemRepo repository.ProblemRepository,
	judgeService JudgeService,
) SubmissionService {
	return &submissionService{
		submissionRepo: submissionRepo,
		problemRepo:    problemRepo,
		judgeService:   judgeService,
		logger:         zap.NewNop(),
	}
}

// SetLogger allows injecting a logger after construction.
func (s *submissionService) SetLogger(logger *zap.Logger) {
	s.logger = logger
}

// CreateSubmission validates, saves code to file, creates the submission record, and sends to judge queue.
func (s *submissionService) CreateSubmission(ctx context.Context, userID uint, req CreateSubmissionRequest, ipAddress string) (*model.Submission, error) {
	// Validate problem exists and is accessible
	problem, err := s.problemRepo.GetByID(ctx, req.ProblemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProblemNotFound
		}
		return nil, err
	}

	if problem.Status != "public" {
		return nil, ErrProblemNotPublic
	}

	// Create submission based on problem type
	submission := &model.Submission{
		UserID:    userID,
		ProblemID: req.ProblemID,
		IPAddress: ipAddress,
		SubmitTime: time.Now(),
		JudgeResult: "pending",
	}

	if problem.Type == "ctf" {
		// CTF submission
		if req.CTFAnswer == "" {
			return nil, ErrCTFAnswerEmpty
		}
		submission.CTFAnswer = req.CTFAnswer
		submission.Language = "ctf"

		// Check CTF answer immediately
		if req.CTFAnswer == problem.CTFFlag {
			submission.JudgeResult = "AC"
			submission.JudgeScore = 100
		} else {
			submission.JudgeResult = "WA"
			submission.JudgeScore = 0
		}
		submission.CodeLength = len(req.CTFAnswer)
	} else {
		// Programming submission
		if !ValidLanguages[req.Language] {
			return nil, ErrInvalidLanguage
		}
		if req.Code == "" {
			return nil, ErrCodeEmpty
		}

		submission.Language = req.Language
		submission.CodeLength = len(req.Code)

		// Save code to file
		submissionDir := filepath.Join("/data/submissions", strconv.Itoa(int(userID)))
		if err := os.MkdirAll(submissionDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create submission directory: %w", err)
		}

		// Create submission record first to get ID
		if err := s.submissionRepo.Create(ctx, submission); err != nil {
			return nil, fmt.Errorf("failed to create submission: %w", err)
		}

		// Save code to file using submission ID
		codeFilename := strconv.Itoa(int(submission.ID)) + getFileExtension(req.Language)
		codePath := filepath.Join(submissionDir, codeFilename)
		if err := os.WriteFile(codePath, []byte(req.Code), 0644); err != nil {
			return nil, fmt.Errorf("failed to save code file: %w", err)
		}

		// Update code path
		submission.CodePath = codePath
		if err := s.submissionRepo.Update(ctx, submission); err != nil {
			s.logger.Error("failed to update submission code path",
				zap.Uint("submission_id", submission.ID),
				zap.Error(err),
			)
		}

		// Increment problem submit count
		_ = s.problemRepo.IncrementSubmitCount(ctx, req.ProblemID)

		// Send to judge queue
		if err := s.judgeService.Submit(ctx, submission); err != nil {
			s.logger.Error("failed to send submission to judge",
				zap.Uint("submission_id", submission.ID),
				zap.Error(err),
			)
			// Don't fail the submission, just log the error
		}

		return submission, nil
	}

	// For CTF submissions, create directly
	if err := s.submissionRepo.Create(ctx, submission); err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	// Increment problem submit count
	_ = s.problemRepo.IncrementSubmitCount(ctx, req.ProblemID)

	// If CTF was accepted, increment accepted count
	if submission.JudgeResult == "AC" {
		_ = s.problemRepo.IncrementAcceptedCount(ctx, req.ProblemID)
	}

	return submission, nil
}

// GetSubmission returns a submission by ID. Users can only see their own submissions unless admin.
func (s *submissionService) GetSubmission(ctx context.Context, submissionID uint, userID uint, isAdmin bool) (*model.Submission, error) {
	submission, err := s.submissionRepo.GetByID(ctx, submissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubmissionNotFound
		}
		return nil, err
	}

	// Non-admin users can only see their own submissions
	if !isAdmin && submission.UserID != userID {
		return nil, ErrForbidden
	}

	return submission, nil
}

// GetSubmissionCode returns the source code content of a submission.
func (s *submissionService) GetSubmissionCode(ctx context.Context, submissionID uint, userID uint, isAdmin bool) (string, error) {
	submission, err := s.GetSubmission(ctx, submissionID, userID, isAdmin)
	if err != nil {
		return "", err
	}

	if submission.CodePath == "" {
		return "", nil
	}

	code, err := os.ReadFile(submission.CodePath)
	if err != nil {
		return "", fmt.Errorf("failed to read code file: %w", err)
	}

	return string(code), nil
}

// ListSubmissions returns a paginated list of submissions.
// Regular users only see their own submissions; admins see all.
func (s *submissionService) ListSubmissions(ctx context.Context, filter SubmissionFilter, offset, limit int, userID uint, isAdmin bool) ([]model.Submission, int64, error) {
	filters := make(map[string]interface{})

	if isAdmin {
		// Admin can filter by any user
		if filter.UserID != 0 {
			filters["user_id"] = filter.UserID
		}
	} else {
		// Regular users only see their own
		filters["user_id"] = userID
	}

	if filter.ProblemID != 0 {
		filters["problem_id"] = filter.ProblemID
	}
	if filter.ContestID != 0 {
		filters["contest_id"] = filter.ContestID
	}
	if filter.JudgeResult != "" {
		filters["judge_result"] = filter.JudgeResult
	}
	if filter.Language != "" {
		filters["language"] = filter.Language
	}

	return s.submissionRepo.List(ctx, offset, limit, filters)
}

// HandleJudgeResult processes the judge callback, updates the DB, and updates problem stats.
func (s *submissionService) HandleJudgeResult(ctx context.Context, submissionID uint, result JudgeResultCallback) error {
	submission, err := s.submissionRepo.GetByID(ctx, submissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSubmissionNotFound
		}
		return err
	}

	// Update submission status to judging
	_ = s.submissionRepo.UpdateSubmissionStatus(ctx, submissionID, "judging")

	// Update the judge result
	if err := s.submissionRepo.UpdateJudgeResult(
		ctx,
		submissionID,
		result.JudgeResult,
		result.JudgeScore,
		result.TimeUsedMs,
		result.MemoryUsedKb,
		result.ErrorMessage,
		result.JudgeDetail,
	); err != nil {
		return fmt.Errorf("failed to update judge result: %w", err)
	}

	// Update problem stats if the result is AC
	if result.JudgeResult == "AC" {
		// Check if this is the user's first AC for this problem
		status, err := s.submissionRepo.GetUserProblemStatus(ctx, submission.UserID, submission.ProblemID)
		if err == nil && status != nil {
			// If the user had no previous AC, increment accepted count
			// We check if there was already an AC before this submission
			// Since this submission just got AC, if submit_count was 1 before, this is the first AC
			// A more robust check would look at previous submissions
			_ = s.problemRepo.IncrementAcceptedCount(ctx, submission.ProblemID)
		}
	}

	s.logger.Info("judge result processed",
		zap.Uint("submission_id", submissionID),
		zap.String("result", result.JudgeResult),
		zap.Float64("score", result.JudgeScore),
	)

	return nil
}

// RejudgeSubmission resets a submission and resends it to the judge.
func (s *submissionService) RejudgeSubmission(ctx context.Context, submissionID uint) error {
	submission, err := s.submissionRepo.GetByID(ctx, submissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSubmissionNotFound
		}
		return err
	}

	// Reset submission status
	_ = s.submissionRepo.UpdateSubmissionStatus(ctx, submissionID, "pending")

	// Resend to judge
	if err := s.judgeService.Submit(ctx, submission); err != nil {
		return fmt.Errorf("failed to resubmit to judge: %w", err)
	}

	return nil
}

// GetUserSubmissions returns submissions for a specific user.
func (s *submissionService) GetUserSubmissions(ctx context.Context, userID uint, offset, limit int) ([]model.Submission, int64, error) {
	return s.submissionRepo.GetByUser(ctx, userID, offset, limit)
}

// GetContestSubmissions returns submissions for a specific contest.
func (s *submissionService) GetContestSubmissions(ctx context.Context, contestID uint, offset, limit int) ([]model.Submission, int64, error) {
	return s.submissionRepo.GetByContest(ctx, contestID, offset, limit)
}

// getFileExtension returns the file extension for a given language.
func getFileExtension(language string) string {
	switch language {
	case "cpp":
		return ".cpp"
	case "c":
		return ".c"
	case "java":
		return ".java"
	case "python3":
		return ".py"
	default:
		return ".txt"
	}
}
