package repository

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"gorm.io/gorm"
)

type ContestRepository interface {
	Create(ctx context.Context, contest *model.Contest) error
	GetByID(ctx context.Context, id uint) (*model.Contest, error)
	Update(ctx context.Context, contest *model.Contest) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Contest, int64, error)
	AddProblem(ctx context.Context, cp *model.ContestProblem) error
	RemoveProblem(ctx context.Context, contestID, problemID uint) error
	GetProblems(ctx context.Context, contestID uint) ([]model.ContestProblem, error)
	Signup(ctx context.Context, signup *model.ContestSignup) error
	Withdraw(ctx context.Context, contestID, userID uint) error
	GetSignup(ctx context.Context, contestID, userID uint) (*model.ContestSignup, error)
	GetSignups(ctx context.Context, contestID uint) ([]model.ContestSignup, error)
	IsSignedUp(ctx context.Context, contestID, userID uint) (bool, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}

type contestRepository struct {
	db *gorm.DB
}

func NewContestRepository(db *gorm.DB) ContestRepository {
	return &contestRepository{db: db}
}

func (r *contestRepository) Create(ctx context.Context, contest *model.Contest) error {
	return r.db.WithContext(ctx).Create(contest).Error
}

func (r *contestRepository) GetByID(ctx context.Context, id uint) (*model.Contest, error) {
	var contest model.Contest
	if err := r.db.WithContext(ctx).Preload("Creator").First(&contest, id).Error; err != nil {
		return nil, err
	}
	return &contest, nil
}

func (r *contestRepository) Update(ctx context.Context, contest *model.Contest) error {
	return r.db.WithContext(ctx).Save(contest).Error
}

func (r *contestRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Contest{}, id).Error
}

func (r *contestRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.Contest, int64, error) {
	var contests []model.Contest
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Contest{})

	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key, value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Creator").Offset(offset).Limit(limit).
		Order("created_at DESC").Find(&contests).Error; err != nil {
		return nil, 0, err
	}

	return contests, total, nil
}

func (r *contestRepository) AddProblem(ctx context.Context, cp *model.ContestProblem) error {
	return r.db.WithContext(ctx).Create(cp).Error
}

func (r *contestRepository) RemoveProblem(ctx context.Context, contestID, problemID uint) error {
	return r.db.WithContext(ctx).Where("contest_id = ? AND problem_id = ?", contestID, problemID).
		Delete(&model.ContestProblem{}).Error
}

func (r *contestRepository) GetProblems(ctx context.Context, contestID uint) ([]model.ContestProblem, error) {
	var problems []model.ContestProblem
	if err := r.db.WithContext(ctx).
		Where("contest_id = ?", contestID).
		Preload("Problem").
		Order("display_order ASC").
		Find(&problems).Error; err != nil {
		return nil, err
	}
	return problems, nil
}

func (r *contestRepository) Signup(ctx context.Context, signup *model.ContestSignup) error {
	return r.db.WithContext(ctx).Create(signup).Error
}

func (r *contestRepository) Withdraw(ctx context.Context, contestID, userID uint) error {
	return r.db.WithContext(ctx).Where("contest_id = ? AND user_id = ?", contestID, userID).
		Delete(&model.ContestSignup{}).Error
}

func (r *contestRepository) GetSignup(ctx context.Context, contestID, userID uint) (*model.ContestSignup, error) {
	var signup model.ContestSignup
	if err := r.db.WithContext(ctx).
		Where("contest_id = ? AND user_id = ?", contestID, userID).
		First(&signup).Error; err != nil {
		return nil, err
	}
	return &signup, nil
}

func (r *contestRepository) GetSignups(ctx context.Context, contestID uint) ([]model.ContestSignup, error) {
	var signups []model.ContestSignup
	if err := r.db.WithContext(ctx).
		Where("contest_id = ?", contestID).
		Preload("User").Preload("Team").
		Find(&signups).Error; err != nil {
		return nil, err
	}
	return signups, nil
}

func (r *contestRepository) IsSignedUp(ctx context.Context, contestID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.ContestSignup{}).
		Where("contest_id = ? AND user_id = ?", contestID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *contestRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&model.Contest{}).Where("id = ?", id).
		Update("status", status).Error
}
