package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
)

type SubmissionService interface {
	CreateSubmission(ctx context.Context, submission *model.Submission) error
	GetSubmission(ctx context.Context, id uint) (*model.Submission, error)
	ListSubmissions(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Submission, int64, error)
	GetUserSubmissions(ctx context.Context, userID uint, offset, limit int) ([]model.Submission, int64, error)
	GetContestSubmissions(ctx context.Context, contestID uint, offset, limit int) ([]model.Submission, int64, error)
}

type submissionService struct {
	submissionRepo repository.SubmissionRepository
}

func NewSubmissionService(submissionRepo repository.SubmissionRepository) SubmissionService {
	return &submissionService{submissionRepo: submissionRepo}
}

func (s *submissionService) CreateSubmission(ctx context.Context, submission *model.Submission) error {
	return s.submissionRepo.Create(ctx, submission)
}

func (s *submissionService) GetSubmission(ctx context.Context, id uint) (*model.Submission, error) {
	return s.submissionRepo.GetByID(ctx, id)
}

func (s *submissionService) ListSubmissions(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Submission, int64, error) {
	return s.submissionRepo.List(ctx, offset, limit, filters)
}

func (s *submissionService) GetUserSubmissions(ctx context.Context, userID uint, offset, limit int) ([]model.Submission, int64, error) {
	return s.submissionRepo.GetByUser(ctx, userID, offset, limit)
}

func (s *submissionService) GetContestSubmissions(ctx context.Context, contestID uint, offset, limit int) ([]model.Submission, int64, error) {
	return s.submissionRepo.GetByContest(ctx, contestID, offset, limit)
}
