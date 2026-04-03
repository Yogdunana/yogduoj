package handler

import (
	"errors"
	"strconv"

	"github.com/Yogdunana/yogduoj/backend/internal/pkg/pagination"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ContestHandler struct {
	contestService service.ContestService
}

func NewContestHandler(contestService service.ContestService) *ContestHandler {
	return &ContestHandler{contestService: contestService}
}

// ListContests returns a paginated list of contests.
// GET /api/v1/contests?page=1&page_size=20&status=running&type=individual&category=programming&rule_type=acm&search=keyword&sort_by=start_time
func (h *ContestHandler) ListContests(c *gin.Context) {
	p := pagination.GetPagination(c)

	filter := service.ContestFilter{
		Status:   c.Query("status"),
		Type:     c.Query("type"),
		Category: c.Query("category"),
		RuleType: c.Query("rule_type"),
		Search:   c.Query("search"),
		SortBy:   c.Query("sort_by"),
	}

	contests, total, err := h.contestService.ListContests(c.Request.Context(), p.Offset(), p.PageSize, filter)
	if err != nil {
		response.InternalError(c, "failed to get contests")
		return
	}

	response.PaginatedResponse(c, contests, total, p.Page, p.PageSize)
}

// GetContest returns a contest by ID with signup status.
// GET /api/v1/contests/:id
func (h *ContestHandler) GetContest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	var userID uint
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	detail, err := h.contestService.GetContestDetail(c.Request.Context(), uint(id), userID)
	if err != nil {
		if errors.Is(err, service.ErrContestNotFound) {
			response.NotFound(c, "contest not found")
			return
		}
		response.InternalError(c, "failed to get contest")
		return
	}

	response.Success(c, detail)
}

// Signup signs up the current user for a contest.
// POST /api/v1/contests/:id/signup
func (h *ContestHandler) Signup(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	if err := h.contestService.SignupUser(c.Request.Context(), uint(contestID), userID.(uint)); err != nil {
		switch {
		case errors.Is(err, service.ErrContestNotFound):
			response.NotFound(c, "contest not found")
		case errors.Is(err, service.ErrContestEnded):
			response.BadRequest(c, "contest has ended")
		case errors.Is(err, service.ErrAlreadySignedUp):
			response.BadRequest(c, "already signed up for this contest")
		case errors.Is(err, service.ErrContestFull):
			response.BadRequest(c, "contest signup limit reached")
		default:
			response.InternalError(c, "failed to sign up for contest")
		}
		return
	}

	response.Success(c, nil)
}

// Withdraw withdraws the current user from a contest.
// POST /api/v1/contests/:id/withdraw
func (h *ContestHandler) Withdraw(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	if err := h.contestService.WithdrawUser(c.Request.Context(), uint(contestID), userID.(uint)); err != nil {
		switch {
		case errors.Is(err, service.ErrContestNotFound):
			response.NotFound(c, "contest not found")
		case errors.Is(err, service.ErrContestNotRunning):
			response.BadRequest(c, "cannot withdraw from running or ended contest")
		default:
			response.InternalError(c, "failed to withdraw from contest")
		}
		return
	}

	response.Success(c, nil)
}

// GetContestProblems returns the problems for a contest.
// GET /api/v1/contests/:id/problems
func (h *ContestHandler) GetContestProblems(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	var userID uint
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	problems, err := h.contestService.GetContestProblems(c.Request.Context(), uint(contestID), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrContestNotFound):
			response.NotFound(c, "contest not found")
		case errors.Is(err, service.ErrNotSignedUp):
			response.Forbidden(c, "not signed up for this contest")
		default:
			response.InternalError(c, "failed to get contest problems")
		}
		return
	}

	response.Success(c, problems)
}

// SubmitToContest creates a submission for a contest problem.
// POST /api/v1/contests/:id/submissions
func (h *ContestHandler) SubmitToContest(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	var req service.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	submission, err := h.contestService.SubmitToContest(
		c.Request.Context(),
		uint(contestID),
		userID.(uint),
		req,
		c.ClientIP(),
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrContestNotFound):
			response.NotFound(c, "contest not found")
		case errors.Is(err, service.ErrContestNotRunning):
			response.BadRequest(c, "contest is not running")
		case errors.Is(err, service.ErrNotSignedUp):
			response.Forbidden(c, "not signed up for this contest")
		case errors.Is(err, service.ErrProblemNotFound):
			response.NotFound(c, "problem not found in this contest")
		case errors.Is(err, service.ErrCodeEmpty):
			response.BadRequest(c, "code cannot be empty")
		case errors.Is(err, service.ErrInvalidLanguage):
			response.BadRequest(c, "invalid programming language")
		case errors.Is(err, service.ErrCTFAnswerEmpty):
			response.BadRequest(c, "ctf answer cannot be empty")
		default:
			response.InternalError(c, "failed to submit to contest")
		}
		return
	}

	response.Success(c, submission)
}

// GetContestRanking returns the ranking for a contest.
// GET /api/v1/contests/:id/ranking
func (h *ContestHandler) GetContestRanking(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	ranking, err := h.contestService.GetRanking(c.Request.Context(), uint(contestID))
	if err != nil {
		if errors.Is(err, service.ErrContestNotFound) {
			response.NotFound(c, "contest not found")
			return
		}
		response.InternalError(c, "failed to get contest ranking")
		return
	}

	response.Success(c, ranking)
}

// GetFrozenRanking returns the frozen ranking for a contest.
// GET /api/v1/contests/:id/ranking/frozen
func (h *ContestHandler) GetFrozenRanking(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	ranking, err := h.contestService.GetFrozenRanking(c.Request.Context(), uint(contestID))
	if err != nil {
		if errors.Is(err, service.ErrContestNotFound) {
			response.NotFound(c, "contest not found")
			return
		}
		response.InternalError(c, "failed to get frozen ranking")
		return
	}

	response.Success(c, ranking)
}
