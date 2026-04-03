package handler

import (
	"strconv"

	"github.com/Yogdunana/yogduoj/backend/internal/pkg/pagination"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMe returns the current user's profile.
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID.(uint))
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, user)
}

// UpdatePassword changes the current user's password.
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	// TODO: implement password change
	response.Error(c, 501, "not implemented")
}

// GetMySubmissions returns the current user's submissions.
func (h *UserHandler) GetMySubmissions(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	p := pagination.GetPagination(c)
	submissions, total, err := h.userService.ListUsers(c.Request.Context(), p.Offset(), p.PageSize)
	if err != nil {
		response.InternalError(c, "failed to get submissions")
		return
	}

	response.PaginatedResponse(c, submissions, total, p.Page, p.PageSize)
}

// GetMyContests returns the current user's contests.
func (h *UserHandler) GetMyContests(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// GetUser returns a user by ID.
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, user)
}
