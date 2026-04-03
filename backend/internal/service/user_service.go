package service

import (
	"context"
	"errors"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/crypto"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/validator"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrUsernameAlreadyModified = errors.New("username can only be modified once")
	ErrPasswordTooShort        = errors.New("new password must be at least 8 characters")
	ErrPasswordWeak            = errors.New("password must contain uppercase, lowercase, and digit")
	ErrSamePassword            = errors.New("new password must be different from old password")
)

type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*model.User, error)
	UpdateProfile(ctx context.Context, userID uint, username, email, avatar string) (*model.User, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserPublicProfile(ctx context.Context, id uint) (*model.User, error)
	ListUsers(ctx context.Context, offset, limit int) ([]model.User, int64, error)
	ListUsersWithFilters(ctx context.Context, offset, limit int, search, role, status string) ([]model.User, int64, error)
	UpdateUser(ctx context.Context, user *model.User) error
	UpdateUserRoleAndStatus(ctx context.Context, userID uint, role, status string) error
	ResetUserPassword(ctx context.Context, userID uint, newPassword string) error
	DeleteUser(ctx context.Context, id uint) error
}

type userService struct {
	userRepo       repository.UserRepository
	submissionRepo repository.SubmissionRepository
	contestRepo    repository.ContestRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// setDependencies allows injecting optional dependencies after construction.
// This is needed because main.go currently creates UserService with only userRepo.
func (s *userService) setDependencies(submissionRepo repository.SubmissionRepository, contestRepo repository.ContestRepository) {
	s.submissionRepo = submissionRepo
	s.contestRepo = contestRepo
}

// NewUserServiceWithDeps creates a UserService with all dependencies.
func NewUserServiceWithDeps(userRepo repository.UserRepository, submissionRepo repository.SubmissionRepository, contestRepo repository.ContestRepository) UserService {
	return &userService{
		userRepo:       userRepo,
		submissionRepo: submissionRepo,
		contestRepo:    contestRepo,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// UpdateProfile updates the user's profile fields.
// Username can only be changed once (if UsernameModified is false).
// Email and avatar can be changed freely.
func (s *userService) UpdateProfile(ctx context.Context, userID uint, username, email, avatar string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Update username if provided and not already modified
	if username != "" && username != user.Username {
		if user.UsernameModified {
			return nil, ErrUsernameAlreadyModified
		}
		if err := validator.ValidateUsername(username); err != nil {
			return nil, err
		}
		// Check if username is already taken
		existing, err := s.userRepo.GetByUsername(ctx, username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil && existing.ID != user.ID {
			return nil, ErrUsernameExists
		}
		user.Username = username
		user.UsernameModified = true
	}

	// Update email if provided
	if email != "" && email != user.Email {
		if err := validator.ValidateEmail(email); err != nil {
			return nil, err
		}
		// Check if email is already taken
		existing, err := s.userRepo.GetByEmail(ctx, email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil && existing.ID != user.ID {
			return nil, ErrEmailExists
		}
		user.Email = email
	}

	// Update avatar if provided
	if avatar != "" {
		user.Avatar = avatar
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword verifies the old password and updates to the new one.
func (s *userService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Verify old password
	if !crypto.CheckPassword(oldPassword, user.PasswordHash) {
		return ErrInvalidPassword
	}

	// Validate new password
	if err := validator.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Ensure new password is different from old
	if crypto.CheckPassword(newPassword, user.PasswordHash) {
		return ErrSamePassword
	}

	// Hash and save new password
	passwordHash, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash

	return s.userRepo.Update(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetUserPublicProfile returns only the public-safe fields of a user.
// Email and phone are cleared before returning.
func (s *userService) GetUserPublicProfile(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	// Clear sensitive fields
	user.Email = ""
	user.Phone = ""
	user.PasswordHash = ""
	return user, nil
}

func (s *userService) ListUsers(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	return s.userRepo.List(ctx, offset, limit)
}

// ListUsersWithFilters returns a paginated list of users with search and filter options.
func (s *userService) ListUsersWithFilters(ctx context.Context, offset, limit int, search, role, status string) ([]model.User, int64, error) {
	return s.userRepo.ListWithFilters(ctx, offset, limit, search, role, status)
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.userRepo.Update(ctx, user)
}

// UpdateUserRoleAndStatus updates only the role and status fields of a user (admin operation).
func (s *userService) UpdateUserRoleAndStatus(ctx context.Context, userID uint, role, status string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if role != "" {
		user.Role = role
	}
	if status != "" {
		user.Status = status
	}

	return s.userRepo.Update(ctx, user)
}

// ResetUserPassword allows an admin to reset a user's password.
func (s *userService) ResetUserPassword(ctx context.Context, userID uint, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if err := validator.ValidatePassword(newPassword); err != nil {
		return err
	}

	passwordHash, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash

	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}
