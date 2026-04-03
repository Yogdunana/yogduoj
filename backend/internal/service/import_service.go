package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
)

// ImportService handles problem importing from external platforms.
type ImportService interface {
	ImportProblems(ctx context.Context, sourcePlatform, filePath string, userID uint) (*model.ImportRecord, error)
	GetImportRecord(ctx context.Context, id uint) (*model.ImportRecord, error)
	ListImportRecords(ctx context.Context, offset, limit int) ([]model.ImportRecord, int64, error)
}

type importService struct {
	// TODO: add repository
}

func NewImportService() ImportService {
	return &importService{}
}

func (s *importService) ImportProblems(ctx context.Context, sourcePlatform, filePath string, userID uint) (*model.ImportRecord, error) {
	// TODO: implement problem import
	return nil, nil
}

func (s *importService) GetImportRecord(ctx context.Context, id uint) (*model.ImportRecord, error) {
	// TODO: implement
	return nil, nil
}

func (s *importService) ListImportRecords(ctx context.Context, offset, limit int) ([]model.ImportRecord, int64, error) {
	// TODO: implement
	return nil, 0, nil
}
