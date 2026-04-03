package handler

import (
	"errors"
	"strconv"

	"github.com/Yogdunana/yogduoj/backend/internal/pkg/pagination"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userService        service.UserService
	problemService     service.ProblemService
	contestService     service.ContestService
	announcementService service.AnnouncementService
	antiCheatService   service.AntiCheatService
	aiService          service.AIService
	importService      service.ImportService
	systemService      service.SystemService
}

func NewAdminHandler(
	userService service.UserService,
	problemService service.ProblemService,
	contestService service.ContestService,
	announcementService service.AnnouncementService,
	antiCheatService service.AntiCheatService,
	aiService service.AIService,
	importService service.ImportService,
	systemService service.SystemService,
) *AdminHandler {
	return &AdminHandler{
		userService:        userService,
		problemService:     problemService,
		contestService:     contestService,
		announcementService: announcementService,
		antiCheatService:   antiCheatService,
		aiService:          aiService,
		importService:      importService,
		systemService:      systemService,
	}
}

// adminUpdateUserRequest is the request body for admin updating a user.
type adminUpdateUserRequest struct {
	Role   string `json:"role" binding:"omitempty,oneof=user admin"`
	Status string `json:"status" binding:"omitempty,oneof=active disabled"`
}

// adminResetPasswordRequest is the request body for admin resetting a user's password.
type adminResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ListUsers returns all users (admin).
// GET /api/v1/admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	p := pagination.GetPagination(c)
	search := c.Query("search")
	role := c.Query("role")
	status := c.Query("status")

	users, total, err := h.userService.ListUsersWithFilters(
		c.Request.Context(),
		p.Offset(),
		p.PageSize,
		search,
		role,
		status,
	)
	if err != nil {
		response.InternalError(c, "failed to get users")
		return
	}

	response.PaginatedResponse(c, users, total, p.Page, p.PageSize)
}

// UpdateUserRole updates a user's role (admin).
// PUT /api/v1/admin/users/:id/role
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required,oneof=user admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	if err := h.userService.UpdateUserRoleAndStatus(c.Request.Context(), uint(userID), req.Role, ""); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, "failed to update user role")
		return
	}

	response.Success(c, nil)
}

// DisableUser disables a user account (admin).
// PUT /api/v1/admin/users/:id/disable
func (h *AdminHandler) DisableUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active disabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	if err := h.userService.UpdateUserRoleAndStatus(c.Request.Context(), uint(userID), "", req.Status); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, "failed to update user status")
		return
	}

	response.Success(c, nil)
}

// UpdateUser updates a user's role and/or status (admin).
// PUT /api/v1/admin/users/:id
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req adminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	if err := h.userService.UpdateUserRoleAndStatus(
		c.Request.Context(),
		uint(userID),
		req.Role,
		req.Status,
	); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, "failed to update user")
		return
	}

	response.Success(c, nil)
}

// ResetUserPassword resets a user's password (admin).
// POST /api/v1/admin/users/:id/reset-password
func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req adminResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	if err := h.userService.ResetUserPassword(c.Request.Context(), uint(userID), req.NewPassword); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		if errors.Is(err, service.ErrPasswordTooShort) {
			response.BadRequest(c, "password must be at least 8 characters")
			return
		}
		if errors.Is(err, service.ErrPasswordWeak) {
			response.BadRequest(c, "password must contain uppercase, lowercase, and digit")
			return
		}
		response.InternalError(c, "failed to reset password")
		return
	}

	response.Success(c, nil)
}

// CreateProblem creates a new problem (admin).
func (h *AdminHandler) CreateProblem(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// UpdateProblem updates a problem (admin).
func (h *AdminHandler) UpdateProblem(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// DeleteProblem deletes a problem (admin).
func (h *AdminHandler) DeleteProblem(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// CreateContest creates a new contest (admin).
func (h *AdminHandler) CreateContest(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// UpdateContest updates a contest (admin).
func (h *AdminHandler) UpdateContest(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// DeleteContest deletes a contest (admin).
func (h *AdminHandler) DeleteContest(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// CreateAnnouncement creates an announcement (admin).
func (h *AdminHandler) CreateAnnouncement(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// DetectCheating triggers cheat detection for a contest (admin).
func (h *AdminHandler) DetectCheating(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// ListCheatRecords returns cheat records (admin).
func (h *AdminHandler) ListCheatRecords(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// ReviewCheatRecord reviews a cheat record (admin).
func (h *AdminHandler) ReviewCheatRecord(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// GenerateAIProblem generates a problem using AI (admin).
func (h *AdminHandler) GenerateAIProblem(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// GenerateAITestdata generates test data using AI (admin).
func (h *AdminHandler) GenerateAITestdata(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// ImportProblems imports problems from external platform (admin).
func (h *AdminHandler) ImportProblems(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// GetSystemConfigs returns system configurations (admin).
func (h *AdminHandler) GetSystemConfigs(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// SetSystemConfig sets a system configuration (admin).
func (h *AdminHandler) SetSystemConfig(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// CreateCTFResource creates a CTF resource (admin).
func (h *AdminHandler) CreateCTFResource(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// ListCTFResources returns CTF resources.
func (h *AdminHandler) ListCTFResources(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}
