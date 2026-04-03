package handler

import (
	"errors"
	"strconv"

	"github.com/Yogdunana/yogduoj/backend/internal/pkg/pagination"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProblemHandler struct {
	problemService service.ProblemService
}

func NewProblemHandler(problemService service.ProblemService) *ProblemHandler {
	return &ProblemHandler{problemService: problemService}
}

// ListProblems returns a paginated list of problems.
// GET /api/v1/problems?page=1&page_size=20&type=programming&difficulty=easy&search=sort&sort=id&order=desc
func (h *ProblemHandler) ListProblems(c *gin.Context) {
	p := pagination.GetPagination(c)

	filter := service.ProblemFilter{
		Type:       c.Query("type"),
		Difficulty: c.Query("difficulty"),
		Status:     c.DefaultQuery("status", "public"),
		Source:     c.Query("source"),
		Search:     c.Query("search"),
		Sort:       c.DefaultQuery("sort", "id"),
		Order:      c.DefaultQuery("order", "desc"),
	}

	// Get optional user ID from context (for authenticated requests)
	var userID uint
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	problems, total, err := h.problemService.ListProblems(c.Request.Context(), filter, p.Offset(), p.PageSize, userID)
	if err != nil {
		response.InternalError(c, "failed to get problems")
		return
	}

	response.PaginatedResponse(c, problems, total, p.Page, p.PageSize)
}

// GetProblem returns a problem by ID with user status.
// GET /api/v1/problems/:id
func (h *ProblemHandler) GetProblem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	// Get optional user ID from context
	var userID uint
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	detail, err := h.problemService.GetProblemDetailResponse(c.Request.Context(), uint(id), userID)
	if err != nil {
		if errors.Is(err, service.ErrProblemNotFound) {
			response.NotFound(c, "problem not found")
			return
		}
		response.InternalError(c, "failed to get problem")
		return
	}

	response.Success(c, detail)
}

// GetProblemSamples returns the samples for a problem.
// GET /api/v1/problems/:id/samples
func (h *ProblemHandler) GetProblemSamples(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	samples, err := h.problemService.GetProblemSamples(c.Request.Context(), uint(id))
	if err != nil {
		response.InternalError(c, "failed to get samples")
		return
	}

	response.Success(c, samples)
}

// GetAttachment returns a problem attachment file.
// GET /api/v1/problems/:id/attachments/:fileId
func (h *ProblemHandler) GetAttachment(c *gin.Context) {
	idStr := c.Param("id")
	problemID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	fileIDStr := c.Param("fileId")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	attachment, err := h.problemService.DownloadAttachment(c.Request.Context(), uint(problemID), uint(fileID))
	if err != nil {
		if errors.Is(err, service.ErrProblemNotFound) {
			response.NotFound(c, "problem not found")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(c, "attachment not found")
			return
		}
		response.InternalError(c, "failed to get attachment")
		return
	}

	c.FileAttachment(attachment.Path, attachment.Name)
}
