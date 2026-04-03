package handler

import (
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
func (h *ProblemHandler) ListProblems(c *gin.Context) {
	p := pagination.GetPagination(c)
	filters := make(map[string]interface{})

	if typ := c.Query("type"); typ != "" {
		filters["type"] = typ
	}
	if difficulty := c.Query("difficulty"); difficulty != "" {
		filters["difficulty"] = difficulty
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if ctfCategory := c.Query("ctf_category"); ctfCategory != "" {
		filters["ctf_category"] = ctfCategory
	}

	problems, total, err := h.problemService.ListProblems(c.Request.Context(), p.Offset(), p.PageSize, filters)
	if err != nil {
		response.InternalError(c, "failed to get problems")
		return
	}

	response.PaginatedResponse(c, problems, total, p.Page, p.PageSize)
}

// GetProblem returns a problem by ID.
func (h *ProblemHandler) GetProblem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	problem, err := h.problemService.GetProblem(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "problem not found")
		return
	}

	response.Success(c, problem)
}

// GetProblemSamples returns the samples for a problem.
func (h *ProblemHandler) GetProblemSamples(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid problem id")
		return
	}

	samples, err := h.problemService.GetSamples(c.Request.Context(), uint(id))
	if err != nil {
		response.InternalError(c, "failed to get samples")
		return
	}

	response.Success(c, samples)
}

// CreateProblem creates a new problem (admin only).
func (h *ProblemHandler) CreateProblem(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// UpdateProblem updates a problem (admin only).
func (h *ProblemHandler) UpdateProblem(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// DeleteProblem deletes a problem (admin only).
func (h *ProblemHandler) DeleteProblem(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// GetAttachment returns a problem attachment file.
func (h *ProblemHandler) GetAttachment(c *gin.Context) {
	// TODO: implement file serving
	response.Error(c, 501, "not implemented")
}
