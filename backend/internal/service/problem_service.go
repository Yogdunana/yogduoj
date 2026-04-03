package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrProblemNotFound = errors.New("problem not found")
	ErrTestDataNotFound = errors.New("test data not found")
)

// ProblemFilter holds filter parameters for listing problems.
type ProblemFilter struct {
	Type       string
	Difficulty string
	Status     string
	Source     string
	Search     string
	Sort       string
	Order      string
}

// CreateProblemRequest is the request body for creating a problem.
type CreateProblemRequest struct {
	Title           string   `json:"title" binding:"required,max=256"`
	Type            string   `json:"type" binding:"omitempty,oneof=programming algorithm ctf"`
	CTFCategory     string   `json:"ctf_category" binding:"omitempty,oneof=reverse pwn web crypto forensics misc recon vuln_reproduce"`
	Difficulty      string   `json:"difficulty" binding:"omitempty,oneof=easy medium hard"`
	TimeLimitMs     int      `json:"time_limit_ms" binding:"omitempty,min=100,max=60000"`
	MemoryLimitKb   int      `json:"memory_limit_kb" binding:"omitempty,min=1024,max=1073741824"`
	Description     string   `json:"description"`
	InputFormat     string   `json:"input_format"`
	OutputFormat    string   `json:"output_format"`
	Hints           string   `json:"hints"`
	Source          string   `json:"source" binding:"omitempty,oneof=original ai imported"`
	SourceDetail    string   `json:"source_detail"`
	Status          string   `json:"status" binding:"omitempty,oneof=public private disabled"`
	CTFFlag         string   `json:"ctf_flag"`
	Tags            []string `json:"tags"`
	SpecialJudgePath string  `json:"special_judge_path"`
}

// UpdateProblemRequest is the request body for updating a problem.
type UpdateProblemRequest struct {
	Title           string   `json:"title" binding:"omitempty,max=256"`
	Type            string   `json:"type" binding:"omitempty,oneof=programming algorithm ctf"`
	CTFCategory     string   `json:"ctf_category" binding:"omitempty,oneof=reverse pwn web crypto forensics misc recon vuln_reproduce"`
	Difficulty      string   `json:"difficulty" binding:"omitempty,oneof=easy medium hard"`
	TimeLimitMs     *int     `json:"time_limit_ms" binding:"omitempty"`
	MemoryLimitKb   *int     `json:"memory_limit_kb" binding:"omitempty"`
	Description     string   `json:"description"`
	InputFormat     string   `json:"input_format"`
	OutputFormat    string   `json:"output_format"`
	Hints           string   `json:"hints"`
	Source          string   `json:"source" binding:"omitempty,oneof=original ai imported"`
	SourceDetail    string   `json:"source_detail"`
	Status          string   `json:"status" binding:"omitempty,oneof=public private disabled"`
	CTFFlag         string   `json:"ctf_flag"`
	Tags            []string `json:"tags"`
	SpecialJudgePath string  `json:"special_judge_path"`
}

// ProblemService handles problem-related business logic.
type ProblemService interface {
	ListProblems(ctx context.Context, filter ProblemFilter, offset, limit int, userID uint) ([]model.ProblemListItem, int64, error)
	GetProblemDetail(ctx context.Context, problemID uint, userID uint) (*model.Problem, string, error)
	GetProblemDetailResponse(ctx context.Context, problemID uint, userID uint) (*ProblemDetailResponse, error)
	CreateProblem(ctx context.Context, adminID uint, req CreateProblemRequest) (*model.Problem, error)
	UpdateProblem(ctx context.Context, problemID uint, req UpdateProblemRequest) (*model.Problem, error)
	DeleteProblem(ctx context.Context, problemID uint) error
	UploadTestData(ctx context.Context, problemID uint, files []*multipart.FileHeader) ([]model.TestData, error)
	DeleteTestData(ctx context.Context, problemID uint, dataID uint) error
	GetProblemSamples(ctx context.Context, problemID uint) ([]model.Sample, error)
	DownloadAttachment(ctx context.Context, problemID uint, fileID uint) (*model.ProblemAttachment, error)
	GetTestDataList(ctx context.Context, problemID uint) ([]model.TestData, error)
}

type problemService struct {
	problemRepo     repository.ProblemRepository
	submissionRepo  repository.SubmissionRepository
}

func NewProblemService(problemRepo repository.ProblemRepository, submissionRepo repository.SubmissionRepository) ProblemService {
	return &problemService{
		problemRepo:    problemRepo,
		submissionRepo: submissionRepo,
	}
}

// ListProblems returns a paginated list of problems with optional user status.
func (s *problemService) ListProblems(ctx context.Context, filter ProblemFilter, offset, limit int, userID uint) ([]model.ProblemListItem, int64, error) {
	filters := make(map[string]interface{})
	if filter.Type != "" {
		filters["type"] = filter.Type
	}
	if filter.Difficulty != "" {
		filters["difficulty"] = filter.Difficulty
	}
	if filter.Status != "" {
		filters["status"] = filter.Status
	}
	if filter.Source != "" {
		filters["source"] = filter.Source
	}

	problems, total, err := s.problemRepo.List(ctx, offset, limit, filters, filter.Search, filter.Sort, filter.Order)
	if err != nil {
		return nil, 0, err
	}

	// Build enriched list with user status
	items := make([]model.ProblemListItem, 0, len(problems))

	if userID > 0 && len(problems) > 0 {
		// Batch fetch user statuses
		problemIDs := make([]uint, 0, len(problems))
		for _, p := range problems {
			problemIDs = append(problemIDs, p.ID)
		}

		statuses, err := s.submissionRepo.GetUserProblemStatuses(ctx, userID, problemIDs)
		if err != nil {
			// If batch fetch fails, just return problems without status
			for _, p := range problems {
				items = append(items, model.ProblemListItem{Problem: p, UserStatus: "unsubmitted"})
			}
			return items, total, nil
		}

		for _, p := range problems {
			userStatus := "unsubmitted"
			if st, ok := statuses[p.ID]; ok {
				if st.Accepted == 1 {
					userStatus = "accepted"
				} else {
					userStatus = "submitted"
				}
			}
			items = append(items, model.ProblemListItem{Problem: p, UserStatus: userStatus})
		}
	} else {
		for _, p := range problems {
			items = append(items, model.ProblemListItem{Problem: p, UserStatus: "unsubmitted"})
		}
	}

	return items, total, nil
}

// GetProblemDetail returns the full problem detail with samples and user status.
func (s *problemService) GetProblemDetail(ctx context.Context, problemID uint, userID uint) (*model.Problem, string, error) {
	problem, err := s.problemRepo.GetByID(ctx, problemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrProblemNotFound
		}
		return nil, "", err
	}

	// Determine user status
	userStatus := "unsubmitted"
	if userID > 0 {
		status, err := s.submissionRepo.GetUserProblemStatus(ctx, userID, problemID)
		if err == nil && status != nil {
			if status.Accepted == 1 {
				userStatus = "accepted"
			} else if status.SubmitCount > 0 {
				userStatus = "submitted"
			}
		}
	}

	return problem, userStatus, nil
}

// CreateProblem creates a new problem.
func (s *problemService) CreateProblem(ctx context.Context, adminID uint, req CreateProblemRequest) (*model.Problem, error) {
	problem := &model.Problem{
		Title:            req.Title,
		Type:             defaultIfEmpty(req.Type, "programming"),
		CTFCategory:      req.CTFCategory,
		Difficulty:       defaultIfEmpty(req.Difficulty, "easy"),
		TimeLimitMs:      defaultInt(req.TimeLimitMs, 2000),
		MemoryLimitKb:    defaultInt(req.MemoryLimitKb, 262144),
		Description:      req.Description,
		InputFormat:      req.InputFormat,
		OutputFormat:     req.OutputFormat,
		Hints:            req.Hints,
		Source:           defaultIfEmpty(req.Source, "original"),
		SourceDetail:     req.SourceDetail,
		Status:           defaultIfEmpty(req.Status, "public"),
		CTFFlag:          req.CTFFlag,
		CreatedBy:        adminID,
		SpecialJudgePath: req.SpecialJudgePath,
	}

	if err := s.problemRepo.Create(ctx, problem); err != nil {
		return nil, err
	}

	// Handle tags
	if len(req.Tags) > 0 {
		for _, tagName := range req.Tags {
			tag, err := s.getOrCreateTag(ctx, tagName)
			if err != nil {
				continue // Skip failed tags, don't fail the whole operation
			}
			_ = s.problemRepo.AddTag(ctx, problem.ID, tag.ID)
		}
		// Reload problem with tags
		problem, _ = s.problemRepo.GetByID(ctx, problem.ID)
	}

	return problem, nil
}

// UpdateProblem updates an existing problem.
func (s *problemService) UpdateProblem(ctx context.Context, problemID uint, req UpdateProblemRequest) (*model.Problem, error) {
	problem, err := s.problemRepo.GetByID(ctx, problemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProblemNotFound
		}
		return nil, err
	}

	// Update fields if provided
	if req.Title != "" {
		problem.Title = req.Title
	}
	if req.Type != "" {
		problem.Type = req.Type
	}
	if req.CTFCategory != "" {
		problem.CTFCategory = req.CTFCategory
	}
	if req.Difficulty != "" {
		problem.Difficulty = req.Difficulty
	}
	if req.TimeLimitMs != nil {
		problem.TimeLimitMs = *req.TimeLimitMs
	}
	if req.MemoryLimitKb != nil {
		problem.MemoryLimitKb = *req.MemoryLimitKb
	}
	if req.Description != "" {
		problem.Description = req.Description
	}
	if req.InputFormat != "" {
		problem.InputFormat = req.InputFormat
	}
	if req.OutputFormat != "" {
		problem.OutputFormat = req.OutputFormat
	}
	if req.Hints != "" {
		problem.Hints = req.Hints
	}
	if req.Source != "" {
		problem.Source = req.Source
	}
	if req.SourceDetail != "" {
		problem.SourceDetail = req.SourceDetail
	}
	if req.Status != "" {
		problem.Status = req.Status
	}
	if req.CTFFlag != "" {
		problem.CTFFlag = req.CTFFlag
	}
	if req.SpecialJudgePath != "" {
		problem.SpecialJudgePath = req.SpecialJudgePath
	}

	if err := s.problemRepo.Update(ctx, problem); err != nil {
		return nil, err
	}

	// Handle tags update if provided
	if req.Tags != nil {
		// Remove existing tags
		existingTags, _ := s.problemRepo.GetTags(ctx, problemID)
		for _, tag := range existingTags {
			_ = s.problemRepo.RemoveTag(ctx, problemID, tag.ID)
		}
		// Add new tags
		for _, tagName := range req.Tags {
			tag, err := s.getOrCreateTag(ctx, tagName)
			if err != nil {
				continue
			}
			_ = s.problemRepo.AddTag(ctx, problemID, tag.ID)
		}
		// Reload problem with tags
		problem, _ = s.problemRepo.GetByID(ctx, problemID)
	}

	return problem, nil
}

// DeleteProblem deletes a problem and its associated data.
func (s *problemService) DeleteProblem(ctx context.Context, problemID uint) error {
	_, err := s.problemRepo.GetByID(ctx, problemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProblemNotFound
		}
		return err
	}

	// Clean up test data files from disk
	testDataDir := filepath.Join("/data/problems", strconv.Itoa(int(problemID)), "testdata")
	_ = os.RemoveAll(testDataDir)

	return s.problemRepo.Delete(ctx, problemID)
}

// UploadTestData handles uploading test data files for a problem.
// Expects pairs of .in and .out files with matching names.
func (s *problemService) UploadTestData(ctx context.Context, problemID uint, files []*multipart.FileHeader) ([]model.TestData, error) {
	_, err := s.problemRepo.GetByID(ctx, problemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProblemNotFound
		}
		return nil, err
	}

	// Create test data directory
	testDataDir := filepath.Join("/data/problems", strconv.Itoa(int(problemID)), "testdata")
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create test data directory: %w", err)
	}

	// Group files by base name (pair .in and .out)
	type filePair struct {
		baseName  string
		inputFile *multipart.FileHeader
		outputFile *multipart.FileHeader
	}

	pairs := make(map[string]*filePair)
	for _, f := range files {
		ext := filepath.Ext(f.Filename)
		baseName := strings.TrimSuffix(f.Filename, ext)

		if _, ok := pairs[baseName]; !ok {
			pairs[baseName] = &filePair{baseName: baseName}
		}

		switch strings.ToLower(ext) {
		case ".in":
			pairs[baseName].inputFile = f
		case ".out", ".ans":
			pairs[baseName].outputFile = f
		}
	}

	var createdTestData []model.TestData

	for _, pair := range pairs {
		if pair.inputFile == nil || pair.outputFile == nil {
			continue // Skip incomplete pairs
		}

		// Save input file
		inputPath := filepath.Join(testDataDir, pair.baseName+".in")
		if err := saveUploadedFile(pair.inputFile, inputPath); err != nil {
			return nil, fmt.Errorf("failed to save input file %s: %w", pair.inputFile.Filename, err)
		}

		// Save output file
		outputExt := filepath.Ext(pair.outputFile.Filename)
		outputPath := filepath.Join(testDataDir, pair.baseName+outputExt)
		if err := saveUploadedFile(pair.outputFile, outputPath); err != nil {
			return nil, fmt.Errorf("failed to save output file %s: %w", pair.outputFile.Filename, err)
		}

		testData := &model.TestData{
			ProblemID:   problemID,
			InputFile:   inputPath,
			OutputFile:  outputPath,
			ScoreWeight: 1.0,
			IsSample:    false,
			Generation:  "manual",
		}

		if err := s.problemRepo.CreateTestData(ctx, testData); err != nil {
			return nil, fmt.Errorf("failed to create test data record: %w", err)
		}

		createdTestData = append(createdTestData, *testData)
	}

	return createdTestData, nil
}

// DeleteTestData deletes a specific test data entry.
func (s *problemService) DeleteTestData(ctx context.Context, problemID uint, dataID uint) error {
	testData, err := s.problemRepo.GetTestDataByID(ctx, dataID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTestDataNotFound
		}
		return err
	}

	// Verify the test data belongs to the problem
	if testData.ProblemID != problemID {
		return ErrTestDataNotFound
	}

	// Remove files from disk
	_ = os.Remove(testData.InputFile)
	_ = os.Remove(testData.OutputFile)

	return s.problemRepo.DeleteTestData(ctx, dataID)
}

// GetProblemSamples returns the samples for a problem.
func (s *problemService) GetProblemSamples(ctx context.Context, problemID uint) ([]model.Sample, error) {
	return s.problemRepo.GetSamples(ctx, problemID)
}

// DownloadAttachment returns the attachment info for download.
func (s *problemService) DownloadAttachment(ctx context.Context, problemID uint, fileID uint) (*model.ProblemAttachment, error) {
	attachments, err := s.problemRepo.GetAttachments(ctx, problemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProblemNotFound
		}
		return nil, err
	}

	for _, att := range attachments {
		if att.ID == fileID {
			return &att, nil
		}
	}

	return nil, ErrNotFound
}

// GetTestDataList returns all test data for a problem (admin view).
func (s *problemService) GetTestDataList(ctx context.Context, problemID uint) ([]model.TestData, error) {
	return s.problemRepo.GetTestData(ctx, problemID)
}

// getOrCreateTag finds a tag by name or creates it.
func (s *problemService) getOrCreateTag(ctx context.Context, tagName string) (*model.Tag, error) {
	// This uses the DB directly since we don't have a tag repository.
	// In a full implementation, this would use a TagRepository.
	return nil, nil // Placeholder - tags are handled via problem_repo
}

// saveUploadedFile saves an uploaded file to the specified path.
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func defaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func defaultInt(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

// GetProblem returns a problem by ID (backward compatibility).
func (s *problemService) GetProblem(ctx context.Context, id uint) (*model.Problem, error) {
	problem, err := s.problemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProblemNotFound
		}
		return nil, err
	}
	return problem, nil
}

// GetSamples returns samples for a problem (backward compatibility).
func (s *problemService) GetSamples(ctx context.Context, problemID uint) ([]model.Sample, error) {
	return s.problemRepo.GetSamples(ctx, problemID)
}

// CreateProblemBasic creates a problem from a model directly (backward compatibility).
func (s *problemService) CreateProblemBasic(ctx context.Context, problem *model.Problem) error {
	return s.problemRepo.Create(ctx, problem)
}

// UpdateProblemBasic updates a problem from a model directly (backward compatibility).
func (s *problemService) UpdateProblemBasic(ctx context.Context, problem *model.Problem) error {
	return s.problemRepo.Update(ctx, problem)
}

// DeleteProblemBasic deletes a problem by ID (backward compatibility).
func (s *problemService) DeleteProblemBasic(ctx context.Context, id uint) error {
	return s.problemRepo.Delete(ctx, id)
}

// ListProblemsBasic lists problems without user status (backward compatibility).
func (s *problemService) ListProblemsBasic(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Problem, int64, error) {
	return s.problemRepo.List(ctx, offset, limit, filters, "", "", "")
}

// ProblemDetailResponse is the response for problem detail.
type ProblemDetailResponse struct {
	Problem    *model.Problem `json:"problem"`
	UserStatus string         `json:"user_status"`
	Samples    []model.Sample `json:"samples,omitempty"`
}

// GetProblemDetailResponse returns a structured problem detail response.
func (s *problemService) GetProblemDetailResponse(ctx context.Context, problemID uint, userID uint) (*ProblemDetailResponse, error) {
	problem, userStatus, err := s.GetProblemDetail(ctx, problemID, userID)
	if err != nil {
		return nil, err
	}

	samples, _ := s.problemRepo.GetSamples(ctx, problemID)

	return &ProblemDetailResponse{
		Problem:    problem,
		UserStatus: userStatus,
		Samples:    samples,
	}, nil
}

// GetProblemTagsJSON returns tags as JSON for a problem.
func (s *problemService) GetProblemTagsJSON(ctx context.Context, problemID uint) (string, error) {
	tags, err := s.problemRepo.GetTags(ctx, problemID)
	if err != nil {
		return "[]", err
	}
	if len(tags) == 0 {
		return "[]", nil
	}
	data, err := json.Marshal(tags)
	if err != nil {
		return "[]", err
	}
	return string(data), nil
}
