package handler

import (
	"net/http"

	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type authResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Register handles user registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrUsernameExists:
			response.Error(c, http.StatusConflict, "username already exists")
		case service.ErrEmailExists:
			response.Error(c, http.StatusConflict, "email already exists")
		default:
			response.InternalError(c, "failed to register user")
		}
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Login handles user login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	accessToken, refreshToken, err := h.authService.Login(
		c.Request.Context(),
		req.Username,
		req.Password,
		c.ClientIP(),
	)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.Unauthorized(c, "user not found")
		case service.ErrInvalidPassword:
			response.Unauthorized(c, "invalid password")
		case service.ErrAccountDisabled:
			response.Forbidden(c, "account is disabled")
		default:
			response.InternalError(c, "login failed")
		}
		return
	}

	response.Success(c, authResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout handles user logout.
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID.(uint)); err != nil {
		response.InternalError(c, "logout failed")
		return
	}

	response.Success(c, nil)
}

// RefreshToken handles token refresh.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters")
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case service.ErrTokenExpired:
			response.Error(c, http.StatusUnauthorized, "refresh token has expired")
		case service.ErrInvalidToken:
			response.Unauthorized(c, "invalid refresh token")
		default:
			response.InternalError(c, "token refresh failed")
		}
		return
	}

	response.Success(c, authResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
