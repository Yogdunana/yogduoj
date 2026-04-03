package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
)

// AntiCheatService handles cheat detection and review.
type AntiCheatService interface {
	DetectCheating(ctx context.Context, contestID uint) ([]model.CheatRecord, error)
	GetCheatRecord(ctx context.Context, id uint) (*model.CheatRecord, error)
	ReviewCheatRecord(ctx context.Context, id uint, status string, penalty string, reviewerID uint) error
	ListCheatRecords(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.CheatRecord, int64, error)
}

type antiCheatService struct {
	// TODO: add repository
}

func NewAntiCheatService() AntiCheatService {
	return &antiCheatService{}
}

func (s *antiCheatService) DetectCheating(ctx context.Context, contestID uint) ([]model.CheatRecord, error) {
	// TODO: implement cheat detection algorithm
	return nil, nil
}

func (s *antiCheatService) GetCheatRecord(ctx context.Context, id uint) (*model.CheatRecord, error) {
	// TODO: implement
	return nil, nil
}

func (s *antiCheatService) ReviewCheatRecord(ctx context.Context, id uint, status string, penalty string, reviewerID uint) error {
	// TODO: implement
	return nil
}

func (s *antiCheatService) ListCheatRecords(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.CheatRecord, int64, error) {
	// TODO: implement
	return nil, 0, nil
}
