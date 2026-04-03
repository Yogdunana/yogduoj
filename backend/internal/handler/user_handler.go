package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Yogdunana/yogduoj/backend/internal/pkg/pagination"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService       service.UserService
	submissionService service.SubmissionService
	contestService    service.ContestService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// NewUserHandlerWithDeps creates a UserHandler with all dependencies for submission/contest listing.
func NewUserHandlerWithDeps(userService service.UserService, submissionService service.SubmissionService, contestService service.ContestService) *UserHandler {
	return &UserHandler{
		userService:       userService,
		submissionService: submissionService,
		contestService:    contestService,
	}
}

// updateProfileRequest is the request body for updating user profile.
type updateProfileRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Avatar   string `json:"avatar" binding:"omitempty"`
}

// changePasswordRequest is the request body for changing password.
type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// GetMe returns the current user's profile.
// GET /api/v1/users/me
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID.(uint))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, "failed to get user profile")
		return
	}

	response.Success(c, user)
}

// UpdateMe updates the current user's profile.
// PUT /api/v1/users/me
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	user, err := h.userService.UpdateProfile(
		c.Request.Context(),
		userID.(uint),
		req.Username,
		req.Email,
		req.Avatar,
	)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.NotFound(c, "user not found")
		case service.ErrUsernameExists:
			response.Error(c, http.StatusConflict, "username already exists")
		case service.ErrEmailExists:
			response.Error(c, http.StatusConflict, "email already exists")
		case service.ErrUsernameAlreadyModified:
			response.BadRequest(c, "username can only be modified once")
		default:
			if errors.Is(err, service.ErrUsernameAlreadyModified) {
				response.BadRequest(c, err.Error())
				return
			}
			response.InternalError(c, "failed to update profile")
		}
		return
	}

	response.Success(c, user)
}

// UpdatePassword changes the current user's password.
// PUT /api/v1/users/me/password
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	err := h.userService.ChangePassword(
		c.Request.Context(),
		userID.(uint),
		req.OldPassword,
		req.NewPassword,
	)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.NotFound(c, "user not found")
		case service.ErrInvalidPassword:
			response.Unauthorized(c, "old password is incorrect")
		case service.ErrPasswordTooShort:
			response.BadRequest(c, "new password must be at least 8 characters")
		case service.ErrPasswordWeak:
			response.BadRequest(c, "password must contain uppercase, lowercase, and digit")
		case service.ErrSamePassword:
			response.BadRequest(c, "new password must be different from old password")
		default:
			response.InternalError(c, "failed to change password")
		}
		return
	}

	response.Success(c, nil)
}

// GetMySubmissions returns the current user's submissions.
// GET /api/v1/users/me/submissions
func (h *UserHandler) GetMySubmissions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	p := pagination.GetPagination(c)

	if h.submissionService != nil {
		submissions, total, err := h.submissionService.GetUserSubmissions(
			c.Request.Context(), userID.(uint), p.Offset(), p.PageSize,
		)
		if err != nil {
			response.InternalError(c, "failed to get submissions")
			return
		}
		response.PaginatedResponse(c, submissions, total, p.Page, p.PageSize)
		return
	}

	response.Error(c, http.StatusNotImplemented, "not implemented")
}

// GetMyContests returns the current user's contest signups.
// GET /api/v1/users/me/contests
func (h *UserHandler) GetMyContests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	p := pagination.GetPagination(c)

	if h.contestService != nil {
		// Use contest service to get user's contests via signups
		// For now, return an empty paginated response since ContestService
		// does not have a dedicated GetUserContests method.
		_ = userID
		_ = p
		response.PaginatedResponse(c, []interface{}{}, 0, p.Page, p.PageSize)
		return
	}

	response.Error(c, http.StatusNotImplemented, "not implemented")
}

// GetUser returns a user's public profile by ID.
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	user, err := h.userService.GetUserPublicProfile(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, "failed to get user profile")
		return
	}

	response.Success(c, user)
}
