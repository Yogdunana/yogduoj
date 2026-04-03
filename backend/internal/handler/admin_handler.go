package handler

import (
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userService       service.UserService
	problemService    service.ProblemService
	contestService    service.ContestService
	announcementService service.AnnouncementService
	antiCheatService  service.AntiCheatService
	aiService         service.AIService
	importService     service.ImportService
	systemService     service.SystemService
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
		userService:       userService,
		problemService:    problemService,
		contestService:    contestService,
		announcementService: announcementService,
		antiCheatService:  antiCheatService,
		aiService:         aiService,
		importService:     importService,
		systemService:     systemService,
	}
}

// ListUsers returns all users (admin).
func (h *AdminHandler) ListUsers(c *gin.Context) {
	// TODO: implement admin user listing
	response.Error(c, 501, "not implemented")
}

// UpdateUserRole updates a user's role (admin).
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// DisableUser disables a user account (admin).
func (h *AdminHandler) DisableUser(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
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
