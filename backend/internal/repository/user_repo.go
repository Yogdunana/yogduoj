package repository

import (
	"context"
	"time"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]model.User, int64, error)
	ListWithFilters(ctx context.Context, offset, limit int, search, role, status string) ([]model.User, int64, error)
	UpdateLastLogin(ctx context.Context, id uint) error
	GetSubmissionCount(ctx context.Context, userID uint) (int64, error)
	GetSolvedCount(ctx context.Context, userID uint) (int64, error)
	IncrementSolvedCount(ctx context.Context, userID uint) error
	GetLoginAttempts(ctx context.Context, userID uint, since time.Time) (int64, error)
	CreateLoginAttempt(ctx context.Context, attempt *model.LoginAttempt) error
	CleanOldLoginAttempts(ctx context.Context, before time.Time) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListWithFilters returns a paginated list of users with optional search and filter parameters.
// search looks for matches in username or email.
// role filters by user role (user, admin).
// status filters by user status (active, disabled).
func (r *userRepository) ListWithFilters(ctx context.Context, offset, limit int, search, role, status string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("username LIKE ? OR email LIKE ?", like, like)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}

// GetSubmissionCount returns the total number of submissions for a user.
func (r *userRepository) GetSubmissionCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Submission{}).
		Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetSolvedCount returns the number of distinct problems a user has solved (AC).
func (r *userRepository) GetSolvedCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Submission{}).
		Where("user_id = ? AND judge_result = ?", userID, "AC").
		Distinct("problem_id").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// IncrementSolvedCount atomically increments the solved_count column for a user.
func (r *userRepository) IncrementSolvedCount(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).
		Update("solved_count", gorm.Expr("solved_count + 1")).Error
}

// GetLoginAttempts returns the number of login attempts for a user since the given time.
func (r *userRepository) GetLoginAttempts(ctx context.Context, userID uint, since time.Time) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.LoginAttempt{}).
		Where("user_id = ? AND attempt_time > ?", userID, since).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CreateLoginAttempt creates a new login attempt record.
func (r *userRepository) CreateLoginAttempt(ctx context.Context, attempt *model.LoginAttempt) error {
	return r.db.WithContext(ctx).Create(attempt).Error
}

// CleanOldLoginAttempts deletes login attempts older than the given time.
func (r *userRepository) CleanOldLoginAttempts(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).Where("attempt_time < ?", before).
		Delete(&model.LoginAttempt{}).Error
}
