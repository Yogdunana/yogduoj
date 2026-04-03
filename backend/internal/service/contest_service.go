package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
)

type ContestService interface {
	CreateContest(ctx context.Context, contest *model.Contest) error
	GetContest(ctx context.Context, id uint) (*model.Contest, error)
	UpdateContest(ctx context.Context, contest *model.Contest) error
	DeleteContest(ctx context.Context, id uint) error
	ListContests(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Contest, int64, error)
	Signup(ctx context.Context, contestID, userID uint) error
	Withdraw(ctx context.Context, contestID, userID uint) error
	GetContestProblems(ctx context.Context, contestID uint) ([]model.ContestProblem, error)
	GetContestRanking(ctx context.Context, contestID uint) (interface{}, error)
}

type contestService struct {
	contestRepo repository.ContestRepository
}

func NewContestService(contestRepo repository.ContestRepository) ContestService {
	return &contestService{contestRepo: contestRepo}
}

func (s *contestService) CreateContest(ctx context.Context, contest *model.Contest) error {
	return s.contestRepo.Create(ctx, contest)
}

func (s *contestService) GetContest(ctx context.Context, id uint) (*model.Contest, error) {
	return s.contestRepo.GetByID(ctx, id)
}

func (s *contestService) UpdateContest(ctx context.Context, contest *model.Contest) error {
	return s.contestRepo.Update(ctx, contest)
}

func (s *contestService) DeleteContest(ctx context.Context, id uint) error {
	return s.contestRepo.Delete(ctx, id)
}

func (s *contestService) ListContests(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Contest, int64, error) {
	return s.contestRepo.List(ctx, offset, limit, filters)
}

func (s *contestService) Signup(ctx context.Context, contestID, userID uint) error {
	signup := &model.ContestSignup{
		ContestID:  contestID,
		UserID:     userID,
	}
	return s.contestRepo.Signup(ctx, signup)
}

func (s *contestService) Withdraw(ctx context.Context, contestID, userID uint) error {
	return s.contestRepo.Withdraw(ctx, contestID, userID)
}

func (s *contestService) GetContestProblems(ctx context.Context, contestID uint) ([]model.ContestProblem, error) {
	return s.contestRepo.GetProblems(ctx, contestID)
}

func (s *contestService) GetContestRanking(ctx context.Context, contestID uint) (interface{}, error) {
	// TODO: implement ranking calculation
	return nil, nil
}
