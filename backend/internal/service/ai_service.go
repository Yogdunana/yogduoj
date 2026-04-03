package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
)

// AIService handles AI-related operations (problem generation, test data generation).
type AIService interface {
	GenerateProblem(ctx context.Context, params map[string]interface{}) (*model.AIProblemRecord, error)
	GenerateTestData(ctx context.Context, problemID uint, params map[string]interface{}) (*model.AITestdataRecord, error)
	ReviewAIProblem(ctx context.Context, id uint, status string) error
	ReviewAITestdata(ctx context.Context, id uint, status string) error
}

type aiService struct {
	// TODO: add AI client and repositories
}

func NewAIService() AIService {
	return &aiService{}
}

func (s *aiService) GenerateProblem(ctx context.Context, params map[string]interface{}) (*model.AIProblemRecord, error) {
	// TODO: implement AI problem generation
	return nil, nil
}

func (s *aiService) GenerateTestData(ctx context.Context, problemID uint, params map[string]interface{}) (*model.AITestdataRecord, error) {
	// TODO: implement AI test data generation
	return nil, nil
}

func (s *aiService) ReviewAIProblem(ctx context.Context, id uint, status string) error {
	// TODO: implement
	return nil
}

func (s *aiService) ReviewAITestdata(ctx context.Context, id uint, status string) error {
	// TODO: implement
	return nil
}
