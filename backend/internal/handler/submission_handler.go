package handler

import (
	"strconv"

	"github.com/Yogdunana/yogduoj/backend/internal/pkg/pagination"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
	submissionService service.SubmissionService
	judgeService      service.JudgeService
}

func NewSubmissionHandler(submissionService service.SubmissionService, judgeService service.JudgeService) *SubmissionHandler {
	return &SubmissionHandler{
		submissionService: submissionService,
		judgeService:      judgeService,
	}
}

// CreateSubmission creates a new submission.
func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	// TODO: implement submission creation
	response.Error(c, 501, "not implemented")
}

// ListSubmissions returns a paginated list of submissions.
func (h *SubmissionHandler) ListSubmissions(c *gin.Context) {
	p := pagination.GetPagination(c)
	filters := make(map[string]interface{})

	if problemID := c.Query("problem_id"); problemID != "" {
		filters["problem_id"] = problemID
	}
	if contestID := c.Query("contest_id"); contestID != "" {
		filters["contest_id"] = contestID
	}
	if result := c.Query("result"); result != "" {
		filters["judge_result"] = result
	}

	submissions, total, err := h.submissionService.ListSubmissions(c.Request.Context(), p.Offset(), p.PageSize, filters)
	if err != nil {
		response.InternalError(c, "failed to get submissions")
		return
	}

	response.PaginatedResponse(c, submissions, total, p.Page, p.PageSize)
}

// GetSubmission returns a submission by ID.
func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid submission id")
		return
	}

	submission, err := h.submissionService.GetSubmission(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "submission not found")
		return
	}

	response.Success(c, submission)
}

// GetSubmissionCode returns the source code of a submission.
func (h *SubmissionHandler) GetSubmissionCode(c *gin.Context) {
	// TODO: implement code retrieval
	response.Error(c, 501, "not implemented")
}

// JudgeWebSocket handles WebSocket connection for real-time judge status.
func (h *SubmissionHandler) JudgeWebSocket(c *gin.Context) {
	// TODO: implement WebSocket for real-time judge updates
	response.Error(c, 501, "not implemented")
}
