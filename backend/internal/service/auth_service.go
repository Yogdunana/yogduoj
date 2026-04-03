package service

import (
	"context"
	"errors"
	"time"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/crypto"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/jwt"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUsernameExists   = errors.New("username already exists")
	ErrEmailExists      = errors.New("email already exists")
	ErrAccountDisabled  = errors.New("account is disabled")
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token has expired")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("resource not found")
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (*model.User, error)
	Login(ctx context.Context, username, password, ipAddress string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, userID uint) error
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(ctx context.Context, username, email, password string) (*model.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUsernameExists
	}

	// Check if email already exists
	if email != "" {
		existingEmail, err := s.userRepo.GetByEmail(ctx, email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingEmail != nil {
			return nil, ErrEmailExists
		}
	}

	// Hash password
	passwordHash, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "user",
		Status:       "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, username, password, ipAddress string) (string, string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrUserNotFound
		}
		return "", "", err
	}

	if user.Status == "disabled" {
		return "", "", ErrAccountDisabled
	}

	if !crypto.CheckPassword(password, user.PasswordHash) {
		return "", "", ErrInvalidPassword
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", "", err
	}

	// Record login attempt
	_ = recordLoginAttempt(ctx, s.userRepo, user.ID, ipAddress, true)

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.jwtManager.ParseToken(refreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrExpiredToken) {
			return "", "", ErrTokenExpired
		}
		return "", "", ErrInvalidToken
	}

	accessToken, err := s.jwtManager.GenerateToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *authService) Logout(ctx context.Context, userID uint) error {
	// In a stateless JWT setup, logout is typically handled client-side
	// by discarding the token. If token blacklisting is needed, implement here.
	return nil
}

// recordLoginAttempt is a helper to record login attempts.
func recordLoginAttempt(ctx context.Context, repo repository.UserRepository, userID uint, ipAddress string, success bool) error {
	// This would use a separate LoginAttempt repository in a full implementation.
	// For now, it's a no-op placeholder.
	_ = ctx
	_ = repo
	_ = userID
	_ = ipAddress
	_ = success
	_ = time.Now()
	return nil
}
