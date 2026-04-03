package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
)

type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*model.User, error)
	UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) error
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	ListUsers(ctx context.Context, offset, limit int) ([]model.User, int64, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*model.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) error {
	// TODO: implement profile update
	return nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	// TODO: implement password change
	return nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	return s.userRepo.List(ctx, offset, limit)
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}
