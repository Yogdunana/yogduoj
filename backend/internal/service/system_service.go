package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
)

// SystemService handles system configuration and CTF resources.
type SystemService interface {
	GetConfig(ctx context.Context, key string) (*model.SystemConfig, error)
	SetConfig(ctx context.Context, key, value, description string) error
	ListConfigs(ctx context.Context) ([]model.SystemConfig, error)

	// CTF Resources
	CreateCTFResource(ctx context.Context, resource *model.CTFResource) error
	GetCTFResource(ctx context.Context, id uint) (*model.CTFResource, error)
	UpdateCTFResource(ctx context.Context, resource *model.CTFResource) error
	DeleteCTFResource(ctx context.Context, id uint) error
	ListCTFResources(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.CTFResource, int64, error)
}

type systemService struct {
	// TODO: add repository
}

func NewSystemService() SystemService {
	return &systemService{}
}

func (s *systemService) GetConfig(ctx context.Context, key string) (*model.SystemConfig, error) {
	// TODO: implement
	return nil, nil
}

func (s *systemService) SetConfig(ctx context.Context, key, value, description string) error {
	// TODO: implement
	return nil
}

func (s *systemService) ListConfigs(ctx context.Context) ([]model.SystemConfig, error) {
	// TODO: implement
	return nil, nil
}

func (s *systemService) CreateCTFResource(ctx context.Context, resource *model.CTFResource) error {
	// TODO: implement
	return nil
}

func (s *systemService) GetCTFResource(ctx context.Context, id uint) (*model.CTFResource, error) {
	// TODO: implement
	return nil, nil
}

func (s *systemService) UpdateCTFResource(ctx context.Context, resource *model.CTFResource) error {
	// TODO: implement
	return nil
}

func (s *systemService) DeleteCTFResource(ctx context.Context, id uint) error {
	// TODO: implement
	return nil
}

func (s *systemService) ListCTFResources(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.CTFResource, int64, error) {
	// TODO: implement
	return nil, 0, nil
}
