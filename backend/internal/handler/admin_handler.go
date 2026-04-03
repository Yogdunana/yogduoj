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
	submissionService  service.SubmissionService
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
	submissionService service.SubmissionService,
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
		submissionService:  submissionService,
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

// ==================== Problem Management (Admin) ====================

// CreateProblem creates a new problem (admin).
// POST /api/v1/admin/problems
// JSON body: { "title": "...", "type": "programming", "difficulty": "easy", ... }
func (h *AdminHandler) CreateProblem(c *gin.Context) {
	var req service.CreateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	problem, err := h.problemService.CreateProblem(c.Request.Context(), adminID.(uint), req)
	if err != nil {
		response.InternalError(c, "failed to create problem")
		return
	}

	response.Success(c, problem)
}

// UpdateProblem updates a problem (admin).
// PUT /api/v1/admin/problems/:id
func (h *AdminHandler) UpdateProblem(c *gin.Context) {
	idStr := c.Param("id")
	problemID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	var req service.UpdateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	problem, err := h.problemService.UpdateProblem(c.Request.Context(), uint(problemID), req)
	if err != nil {
		if errors.Is(err, service.ErrProblemNotFound) {
			response.NotFound(c, "problem not found")
			return
		}
		response.InternalError(c, "failed to update problem")
		return
	}

	response.Success(c, problem)
}

// DeleteProblem deletes a problem (admin).
// DELETE /api/v1/admin/problems/:id
func (h *AdminHandler) DeleteProblem(c *gin.Context) {
	idStr := c.Param("id")
	problemID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	if err := h.problemService.DeleteProblem(c.Request.Context(), uint(problemID)); err != nil {
		if errors.Is(err, service.ErrProblemNotFound) {
			response.NotFound(c, "problem not found")
			return
		}
		response.InternalError(c, "failed to delete problem")
		return
	}

	response.Success(c, nil)
}

// UploadTestData uploads test data files for a problem (admin).
// POST /api/v1/admin/problems/:id/testdata
// Multipart form: files[] (pairs of .in and .out files)
func (h *AdminHandler) UploadTestData(c *gin.Context) {
	idStr := c.Param("id")
	problemID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, "failed to parse multipart form: "+err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.BadRequest(c, "no files uploaded")
		return
	}

	testData, err := h.problemService.UploadTestData(c.Request.Context(), uint(problemID), files)
	if err != nil {
		if errors.Is(err, service.ErrProblemNotFound) {
			response.NotFound(c, "problem not found")
			return
		}
		response.InternalError(c, "failed to upload test data: "+err.Error())
		return
	}

	response.Success(c, testData)
}

// DeleteTestData deletes a specific test data entry (admin).
// DELETE /api/v1/admin/problems/:id/testdata/:dataId
func (h *AdminHandler) DeleteTestData(c *gin.Context) {
	idStr := c.Param("id")
	problemID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	dataIdStr := c.Param("dataId")
	dataID, err := strconv.ParseUint(dataIdStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid test data id")
		return
	}

	if err := h.problemService.DeleteTestData(c.Request.Context(), uint(problemID), uint(dataID)); err != nil {
		if errors.Is(err, service.ErrTestDataNotFound) {
			response.NotFound(c, "test data not found")
			return
		}
		response.InternalError(c, "failed to delete test data")
		return
	}

	response.Success(c, nil)
}

// GetTestDataList returns all test data for a problem (admin).
// GET /api/v1/admin/problems/:id/testdata
func (h *AdminHandler) GetTestDataList(c *gin.Context) {
	idStr := c.Param("id")
	problemID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	testData, err := h.problemService.GetTestDataList(c.Request.Context(), uint(problemID))
	if err != nil {
		response.InternalError(c, "failed to get test data")
		return
	}

	response.Success(c, testData)
}

// ==================== Submission Management (Admin) ====================

// ListAllSubmissions returns all submissions (admin sees all).
// GET /api/v1/admin/submissions?page=1&page_size=20&problem_id=1&user_id=1&result=AC&language=cpp
func (h *AdminHandler) ListAllSubmissions(c *gin.Context) {
	p := pagination.GetPagination(c)

	filter := service.SubmissionFilter{}

	if problemID := c.Query("problem_id"); problemID != "" {
		if id, err := strconv.ParseUint(problemID, 10, 32); err == nil {
			filter.ProblemID = uint(id)
		}
	}
	if userID := c.Query("user_id"); userID != "" {
		if id, err := strconv.ParseUint(userID, 10, 32); err == nil {
			filter.UserID = uint(id)
		}
	}
	if contestID := c.Query("contest_id"); contestID != "" {
		if id, err := strconv.ParseUint(contestID, 10, 32); err == nil {
			filter.ContestID = uint(id)
		}
	}
	if result := c.Query("result"); result != "" {
		filter.JudgeResult = result
	}
	if language := c.Query("language"); language != "" {
		filter.Language = language
	}

	submissions, total, err := h.submissionService.ListSubmissions(
		c.Request.Context(),
		filter,
		p.Offset(),
		p.PageSize,
		0, // Admin sees all, no user filter
		true, // isAdmin = true
	)
	if err != nil {
		response.InternalError(c, "failed to get submissions")
		return
	}

	response.PaginatedResponse(c, submissions, total, p.Page, p.PageSize)
}

// RejudgeSubmission triggers rejudge for a submission (admin).
// POST /api/v1/admin/submissions/:id/rejudge
func (h *AdminHandler) RejudgeSubmission(c *gin.Context) {
	idStr := c.Param("id")
	submissionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid submission id")
		return
	}

	if err := h.submissionService.RejudgeSubmission(c.Request.Context(), uint(submissionID)); err != nil {
		if errors.Is(err, service.ErrSubmissionNotFound) {
			response.NotFound(c, "submission not found")
			return
		}
		response.InternalError(c, "failed to rejudge submission")
		return
	}

	response.Success(c, gin.H{
		"message": "rejudge triggered",
	})
}

// ==================== Contest Management (Admin) ====================

// CreateContest creates a new contest (admin).
// POST /api/v1/admin/contests
func (h *AdminHandler) CreateContest(c *gin.Context) {
	var req service.CreateContestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	contest, err := h.contestService.CreateContest(c.Request.Context(), adminID.(uint), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidContestTime):
			response.BadRequest(c, "end time must be after start time")
		default:
			response.InternalError(c, "failed to create contest")
		}
		return
	}

	response.Success(c, contest)
}

// UpdateContest updates a contest (admin).
// PUT /api/v1/admin/contests/:id
func (h *AdminHandler) UpdateContest(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	var req service.UpdateContestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	contest, err := h.contestService.UpdateContest(c.Request.Context(), uint(contestID), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrContestNotFound):
			response.NotFound(c, "contest not found")
		case errors.Is(err, service.ErrInvalidContestTime):
			response.BadRequest(c, "end time must be after start time")
		default:
			response.InternalError(c, "failed to update contest")
		}
		return
	}

	response.Success(c, contest)
}

// DeleteContest deletes a contest (admin).
// DELETE /api/v1/admin/contests/:id
func (h *AdminHandler) DeleteContest(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	if err := h.contestService.DeleteContest(c.Request.Context(), uint(contestID)); err != nil {
		if errors.Is(err, service.ErrContestNotFound) {
			response.NotFound(c, "contest not found")
			return
		}
		response.InternalError(c, "failed to delete contest")
		return
	}

	response.Success(c, nil)
}

// UpdateContestStatus updates a contest's status (admin).
// PUT /api/v1/admin/contests/:id/status
func (h *AdminHandler) UpdateContestStatus(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending running ended cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	if err := h.contestService.UpdateContestStatus(c.Request.Context(), uint(contestID), req.Status); err != nil {
		switch {
		case errors.Is(err, service.ErrContestNotFound):
			response.NotFound(c, "contest not found")
		case errors.Is(err, service.ErrInvalidContestStatus):
			response.BadRequest(c, "invalid contest status")
		default:
			response.InternalError(c, "failed to update contest status")
		}
		return
	}

	response.Success(c, nil)
}

// AddContestProblem adds a problem to a contest (admin).
// POST /api/v1/admin/contests/:id/problems
func (h *AdminHandler) AddContestProblem(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	var req service.ContestProblemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	if err := h.contestService.AddContestProblem(c.Request.Context(), uint(contestID), req); err != nil {
		if errors.Is(err, service.ErrContestNotFound) {
			response.NotFound(c, "contest not found")
			return
		}
		response.InternalError(c, "failed to add problem to contest")
		return
	}

	response.Success(c, nil)
}

// RemoveContestProblem removes a problem from a contest (admin).
// DELETE /api/v1/admin/contests/:id/problems/:problemId
func (h *AdminHandler) RemoveContestProblem(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	problemIdStr := c.Param("problemId")
	problemID, err := strconv.ParseUint(problemIdStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	if err := h.contestService.RemoveContestProblem(c.Request.Context(), uint(contestID), uint(problemID)); err != nil {
		if errors.Is(err, service.ErrContestNotFound) {
			response.NotFound(c, "contest not found")
			return
		}
		response.InternalError(c, "failed to remove problem from contest")
		return
	}

	response.Success(c, nil)
}

// GetContestSignups returns all signups for a contest (admin).
// GET /api/v1/admin/contests/:id/signups
func (h *AdminHandler) GetContestSignups(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	signups, err := h.contestService.GetContestSignups(c.Request.Context(), uint(contestID))
	if err != nil {
		if errors.Is(err, service.ErrContestNotFound) {
			response.NotFound(c, "contest not found")
			return
		}
		response.InternalError(c, "failed to get contest signups")
		return
	}

	response.Success(c, signups)
}

// ==================== DIY Template Management (Admin) ====================

// CreateDIYTemplate creates a new DIY contest template (admin).
// POST /api/v1/admin/diy-templates
func (h *AdminHandler) CreateDIYTemplate(c *gin.Context) {
	var req service.CreateDIYTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	tmpl, err := h.contestService.CreateDIYTemplate(c.Request.Context(), adminID.(uint), req)
	if err != nil {
		response.InternalError(c, "failed to create DIY template")
		return
	}

	response.Success(c, tmpl)
}

// ListDIYTemplates returns all DIY templates (admin).
// GET /api/v1/admin/diy-templates
func (h *AdminHandler) ListDIYTemplates(c *gin.Context) {
	p := pagination.GetPagination(c)

	templates, total, err := h.contestService.ListDIYTemplates(c.Request.Context(), p.Offset(), p.PageSize)
	if err != nil {
		response.InternalError(c, "failed to get DIY templates")
		return
	}

	response.PaginatedResponse(c, templates, total, p.Page, p.PageSize)
}

// UpdateDIYTemplate updates a DIY template (admin).
// PUT /api/v1/admin/diy-templates/:id
func (h *AdminHandler) UpdateDIYTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	var req service.UpdateDIYTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	tmpl, err := h.contestService.UpdateDIYTemplate(c.Request.Context(), uint(id), req)
	if err != nil {
		if errors.Is(err, service.ErrContestNotFound) {
			response.NotFound(c, "DIY template not found")
			return
		}
		response.InternalError(c, "failed to update DIY template")
		return
	}

	response.Success(c, tmpl)
}

// DeleteDIYTemplate deletes a DIY template (admin).
// DELETE /api/v1/admin/diy-templates/:id
func (h *AdminHandler) DeleteDIYTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	if err := h.contestService.DeleteDIYTemplate(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrContestNotFound) {
			response.NotFound(c, "DIY template not found")
			return
		}
		response.InternalError(c, "failed to delete DIY template")
		return
	}

	response.Success(c, nil)
}

// ==================== Announcement Management (Admin) ====================

// CreateAnnouncement creates an announcement (admin).
// POST /api/v1/admin/announcements
func (h *AdminHandler) CreateAnnouncement(c *gin.Context) {
	var req service.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	announcement, err := h.announcementService.CreateAnnouncement(c.Request.Context(), adminID.(uint), req)
	if err != nil {
		response.InternalError(c, "failed to create announcement")
		return
	}

	response.Success(c, announcement)
}

// UpdateAnnouncement updates an announcement (admin).
// PUT /api/v1/admin/announcements/:id
func (h *AdminHandler) UpdateAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid announcement id")
		return
	}

	var req service.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	announcement, err := h.announcementService.UpdateAnnouncement(c.Request.Context(), uint(id), req)
	if err != nil {
		if errors.Is(err, service.ErrAnnouncementNotFound) {
			response.NotFound(c, "announcement not found")
			return
		}
		response.InternalError(c, "failed to update announcement")
		return
	}

	response.Success(c, announcement)
}

// DeleteAnnouncement deletes an announcement (admin).
// DELETE /api/v1/admin/announcements/:id
func (h *AdminHandler) DeleteAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid announcement id")
		return
	}

	if err := h.announcementService.DeleteAnnouncement(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrAnnouncementNotFound) {
			response.NotFound(c, "announcement not found")
			return
		}
		response.InternalError(c, "failed to delete announcement")
		return
	}

	response.Success(c, nil)
}

// ==================== Anti-Cheat (Admin) ====================

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

// ==================== AI Features (Admin) ====================

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

// ==================== Import (Admin) ====================

// ImportProblems imports problems from external platform (admin).
func (h *AdminHandler) ImportProblems(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// ==================== System Config (Admin) ====================

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

// ==================== CTF Resources (Admin) ====================

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
