package repository

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	Create(ctx context.Context, submission *model.Submission) error
	GetByID(ctx context.Context, id uint) (*model.Submission, error)
	Update(ctx context.Context, submission *model.Submission) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Submission, int64, error)
	GetByUser(ctx context.Context, userID uint, offset, limit int) ([]model.Submission, int64, error)
	GetByProblem(ctx context.Context, problemID uint, offset, limit int) ([]model.Submission, int64, error)
	GetByContest(ctx context.Context, contestID uint, offset, limit int) ([]model.Submission, int64, error)
	UpdateJudgeResult(ctx context.Context, id uint, result string, score float64, timeUsed, memoryUsed int, errMsg string) error
	CountByProblemAndResult(ctx context.Context, problemID uint, result string) (int64, error)
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) Create(ctx context.Context, submission *model.Submission) error {
	return r.db.WithContext(ctx).Create(submission).Error
}

func (r *submissionRepository) GetByID(ctx context.Context, id uint) (*model.Submission, error) {
	var submission model.Submission
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Problem").
		First(&submission, id).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *submissionRepository) Update(ctx context.Context, submission *model.Submission) error {
	return r.db.WithContext(ctx).Save(submission).Error
}

func (r *submissionRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Submission, int64, error) {
	var submissions []model.Submission
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Submission{})

	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key, value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Preload("Problem").Offset(offset).Limit(limit).
		Order("submit_time DESC").Find(&submissions).Error; err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

func (r *submissionRepository) GetByUser(ctx context.Context, userID uint, offset, limit int) ([]model.Submission, int64, error) {
	var submissions []model.Submission
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Submission{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Problem").Offset(offset).Limit(limit).
		Order("submit_time DESC").Find(&submissions).Error; err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

func (r *submissionRepository) GetByProblem(ctx context.Context, problemID uint, offset, limit int) ([]model.Submission, int64, error) {
	var submissions []model.Submission
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Submission{}).Where("problem_id = ?", problemID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Offset(offset).Limit(limit).
		Order("submit_time DESC").Find(&submissions).Error; err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

func (r *submissionRepository) GetByContest(ctx context.Context, contestID uint, offset, limit int) ([]model.Submission, int64, error) {
	var submissions []model.Submission
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Submission{}).Where("contest_id = ?", contestID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Preload("Problem").Offset(offset).Limit(limit).
		Order("submit_time DESC").Find(&submissions).Error; err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

func (r *submissionRepository) UpdateJudgeResult(ctx context.Context, id uint, result string, score float64, timeUsed, memoryUsed int, errMsg string) error {
	return r.db.WithContext(ctx).Model(&model.Submission{}).Where("id = ?", id).Updates(map[string]interface{}{
		"judge_result":  result,
		"judge_score":   score,
		"time_used_ms":  timeUsed,
		"memory_used_kb": memoryUsed,
		"error_message": errMsg,
		"judge_end":     gorm.Expr("NOW()"),
	}).Error
}

func (r *submissionRepository) CountByProblemAndResult(ctx context.Context, problemID uint, result string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Submission{}).
		Where("problem_id = ? AND judge_result = ?", problemID, result).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
