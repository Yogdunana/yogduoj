package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
)

type ProblemService interface {
	CreateProblem(ctx context.Context, problem *model.Problem) error
	GetProblem(ctx context.Context, id uint) (*model.Problem, error)
	UpdateProblem(ctx context.Context, problem *model.Problem) error
	DeleteProblem(ctx context.Context, id uint) error
	ListProblems(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Problem, int64, error)
	GetSamples(ctx context.Context, problemID uint) ([]model.Sample, error)
}

type problemService struct {
	problemRepo repository.ProblemRepository
}

func NewProblemService(problemRepo repository.ProblemRepository) ProblemService {
	return &problemService{problemRepo: problemRepo}
}

func (s *problemService) CreateProblem(ctx context.Context, problem *model.Problem) error {
	return s.problemRepo.Create(ctx, problem)
}

func (s *problemService) GetProblem(ctx context.Context, id uint) (*model.Problem, error) {
	return s.problemRepo.GetByID(ctx, id)
}

func (s *problemService) UpdateProblem(ctx context.Context, problem *model.Problem) error {
	return s.problemRepo.Update(ctx, problem)
}

func (s *problemService) DeleteProblem(ctx context.Context, id uint) error {
	return s.problemRepo.Delete(ctx, id)
}

func (s *problemService) ListProblems(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Problem, int64, error) {
	return s.problemRepo.List(ctx, offset, limit, filters)
}

func (s *problemService) GetSamples(ctx context.Context, problemID uint) ([]model.Sample, error) {
	return s.problemRepo.GetSamples(ctx, problemID)
}
