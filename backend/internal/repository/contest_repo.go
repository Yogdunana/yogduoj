package repository

import (
	"context"
	"time"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"gorm.io/gorm"
)

type ContestRepository interface {
	Create(ctx context.Context, contest *model.Contest) error
	GetByID(ctx context.Context, id uint) (*model.Contest, error)
	Update(ctx context.Context, contest *model.Contest) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}, search string, sortBy string) ([]model.Contest, int64, error)
	AddProblem(ctx context.Context, cp *model.ContestProblem) error
	RemoveProblem(ctx context.Context, contestID, problemID uint) error
	GetProblems(ctx context.Context, contestID uint) ([]model.ContestProblem, error)
	Signup(ctx context.Context, signup *model.ContestSignup) error
	Withdraw(ctx context.Context, contestID, userID uint) error
	GetSignup(ctx context.Context, contestID, userID uint) (*model.ContestSignup, error)
	GetSignups(ctx context.Context, contestID uint) ([]model.ContestSignup, error)
	GetSignupsByContest(ctx context.Context, contestID uint) ([]model.ContestSignup, error)
	IsSignedUp(ctx context.Context, contestID, userID uint) (bool, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	GetParticipants(ctx context.Context, contestID uint) ([]model.ContestSignup, error)
	GetContestRanking(ctx context.Context, contestID uint) ([]model.Submission, error)
	GetFrozenRanking(ctx context.Context, contestID uint, freezeTime time.Time) ([]model.Submission, error)
	CountSignups(ctx context.Context, contestID uint) (int64, error)
	IncrementParticipantCount(ctx context.Context, contestID uint) error
	DecrementParticipantCount(ctx context.Context, contestID uint) error
	GetContestsNeedStatusUpdate(ctx context.Context) ([]model.Contest, error)
	// DIY template CRUD
	CreateDIYTemplate(ctx context.Context, tmpl *model.DIYContestTemplate) error
	GetDIYTemplateByID(ctx context.Context, id uint) (*model.DIYContestTemplate, error)
	UpdateDIYTemplate(ctx context.Context, tmpl *model.DIYContestTemplate) error
	DeleteDIYTemplate(ctx context.Context, id uint) error
	ListDIYTemplates(ctx context.Context, offset, limit int) ([]model.DIYContestTemplate, int64, error)
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

func (r *contestRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}, search string, sortBy string) ([]model.Contest, int64, error) {
	var contests []model.Contest
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Contest{})

	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key, value)
		}
	}

	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := "created_at DESC"
	switch sortBy {
	case "start_time":
		orderClause = "start_time DESC"
	case "created_at":
		orderClause = "created_at DESC"
	}

	if err := query.Preload("Creator").Offset(offset).Limit(limit).
		Order(orderClause).Find(&contests).Error; err != nil {
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

func (r *contestRepository) GetSignupsByContest(ctx context.Context, contestID uint) ([]model.ContestSignup, error) {
	return r.GetSignups(ctx, contestID)
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

func (r *contestRepository) GetParticipants(ctx context.Context, contestID uint) ([]model.ContestSignup, error) {
	var signups []model.ContestSignup
	if err := r.db.WithContext(ctx).
		Where("contest_id = ?", contestID).
		Preload("User").
		Order("signup_time ASC").
		Find(&signups).Error; err != nil {
		return nil, err
	}
	return signups, nil
}

func (r *contestRepository) GetContestRanking(ctx context.Context, contestID uint) ([]model.Submission, error) {
	var submissions []model.Submission
	if err := r.db.WithContext(ctx).
		Where("contest_id = ?", contestID).
		Preload("User").Preload("Problem").
		Order("submit_time ASC").
		Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *contestRepository) GetFrozenRanking(ctx context.Context, contestID uint, freezeTime time.Time) ([]model.Submission, error) {
	var submissions []model.Submission
	if err := r.db.WithContext(ctx).
		Where("contest_id = ? AND submit_time <= ?", contestID, freezeTime).
		Preload("User").Preload("Problem").
		Order("submit_time ASC").
		Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *contestRepository) CountSignups(ctx context.Context, contestID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.ContestSignup{}).
		Where("contest_id = ?", contestID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *contestRepository) IncrementParticipantCount(ctx context.Context, contestID uint) error {
	return r.db.WithContext(ctx).Model(&model.Contest{}).
		Where("id = ?", contestID).
		UpdateColumn("participant_count", gorm.Expr("participant_count + 1")).Error
}

func (r *contestRepository) DecrementParticipantCount(ctx context.Context, contestID uint) error {
	return r.db.WithContext(ctx).Model(&model.Contest{}).
		Where("id = ? AND participant_count > 0", contestID).
		UpdateColumn("participant_count", gorm.Expr("participant_count - 1")).Error
}

func (r *contestRepository) GetContestsNeedStatusUpdate(ctx context.Context) ([]model.Contest, error) {
	var contests []model.Contest
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Where("(status = ? AND start_time <= ?) OR (status = ? AND end_time <= ?)", "pending", now, "running", now).
		Find(&contests).Error; err != nil {
		return nil, err
	}
	return contests, nil
}

// ==================== DIY Contest Template ====================

func (r *contestRepository) CreateDIYTemplate(ctx context.Context, tmpl *model.DIYContestTemplate) error {
	return r.db.WithContext(ctx).Create(tmpl).Error
}

func (r *contestRepository) GetDIYTemplateByID(ctx context.Context, id uint) (*model.DIYContestTemplate, error) {
	var tmpl model.DIYContestTemplate
	if err := r.db.WithContext(ctx).Preload("Creator").First(&tmpl, id).Error; err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (r *contestRepository) UpdateDIYTemplate(ctx context.Context, tmpl *model.DIYContestTemplate) error {
	return r.db.WithContext(ctx).Save(tmpl).Error
}

func (r *contestRepository) DeleteDIYTemplate(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.DIYContestTemplate{}, id).Error
}

func (r *contestRepository) ListDIYTemplates(ctx context.Context, offset, limit int) ([]model.DIYContestTemplate, int64, error) {
	var templates []model.DIYContestTemplate
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.DIYContestTemplate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Creator").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}
