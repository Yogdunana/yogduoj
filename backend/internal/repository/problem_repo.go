package repository

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"gorm.io/gorm"
)

type ProblemRepository interface {
	Create(ctx context.Context, problem *model.Problem) error
	GetByID(ctx context.Context, id uint) (*model.Problem, error)
	Update(ctx context.Context, problem *model.Problem) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Problem, int64, error)
	GetByContestID(ctx context.Context, contestID uint) ([]model.Problem, error)
	AddTag(ctx context.Context, problemID, tagID uint) error
	RemoveTag(ctx context.Context, problemID, tagID uint) error
	GetTags(ctx context.Context, problemID uint) ([]model.Tag, error)
	CreateTestData(ctx context.Context, testData *model.TestData) error
	GetTestData(ctx context.Context, problemID uint) ([]model.TestData, error)
	GetSamples(ctx context.Context, problemID uint) ([]model.Sample, error)
}

type problemRepository struct {
	db *gorm.DB
}

func NewProblemRepository(db *gorm.DB) ProblemRepository {
	return &problemRepository{db: db}
}

func (r *problemRepository) Create(ctx context.Context, problem *model.Problem) error {
	return r.db.WithContext(ctx).Create(problem).Error
}

func (r *problemRepository) GetByID(ctx context.Context, id uint) (*model.Problem, error) {
	var problem model.Problem
	if err := r.db.WithContext(ctx).Preload("Tags").Preload("Creator").First(&problem, id).Error; err != nil {
		return nil, err
	}
	return &problem, nil
}

func (r *problemRepository) Update(ctx context.Context, problem *model.Problem) error {
	return r.db.WithContext(ctx).Save(problem).Error
}

func (r *problemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Problem{}, id).Error
}

func (r *problemRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Problem, int64, error) {
	var problems []model.Problem
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Problem{})

	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key, value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Tags").Offset(offset).Limit(limit).Find(&problems).Error; err != nil {
		return nil, 0, err
	}

	return problems, total, nil
}

func (r *problemRepository) GetByContestID(ctx context.Context, contestID uint) ([]model.Problem, error) {
	var problems []model.Problem
	if err := r.db.WithContext(ctx).
		Joins("JOIN contest_problems ON contest_problems.problem_id = problems.id").
		Where("contest_problems.contest_id = ?", contestID).
		Find(&problems).Error; err != nil {
		return nil, err
	}
	return problems, nil
}

func (r *problemRepository) AddTag(ctx context.Context, problemID, tagID uint) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT IGNORE INTO problem_tags (problem_id, tag_id) VALUES (?, ?)",
		problemID, tagID,
	).Error
}

func (r *problemRepository) RemoveTag(ctx context.Context, problemID, tagID uint) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM problem_tags WHERE problem_id = ? AND tag_id = ?",
		problemID, tagID,
	).Error
}

func (r *problemRepository) GetTags(ctx context.Context, problemID uint) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.db.WithContext(ctx).
		Joins("JOIN problem_tags ON problem_tags.tag_id = tags.id").
		Where("problem_tags.problem_id = ?", problemID).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *problemRepository) CreateTestData(ctx context.Context, testData *model.TestData) error {
	return r.db.WithContext(ctx).Create(testData).Error
}

func (r *problemRepository) GetTestData(ctx context.Context, problemID uint) ([]model.TestData, error) {
	var testData []model.TestData
	if err := r.db.WithContext(ctx).Where("problem_id = ?", problemID).Find(&testData).Error; err != nil {
		return nil, err
	}
	return testData, nil
}

func (r *problemRepository) GetSamples(ctx context.Context, problemID uint) ([]model.Sample, error) {
	var samples []model.Sample
	if err := r.db.WithContext(ctx).
		Where("problem_id = ?", problemID).
		Order("display_order ASC").
		Preload("TestData").
		Find(&samples).Error; err != nil {
		return nil, err
	}
	return samples, nil
}
