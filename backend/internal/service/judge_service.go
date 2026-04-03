package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
)

// JudgeService handles code judging operations.
type JudgeService interface {
	Submit(ctx context.Context, submission *model.Submission) error
	GetJudgeStatus(ctx context.Context, submissionID uint) (*model.Submission, error)
}

type judgeService struct {
	// TODO: add judge client (e.g., gRPC connection to judge server)
}

func NewJudgeService() JudgeService {
	return &judgeService{}
}

func (s *judgeService) Submit(ctx context.Context, submission *model.Submission) error {
	// TODO: implement judge submission - send to judge server
	submission.JudgeResult = "pending"
	return nil
}

func (s *judgeService) GetJudgeStatus(ctx context.Context, submissionID uint) (*model.Submission, error) {
	// TODO: implement judge status retrieval
	return nil, nil
}
