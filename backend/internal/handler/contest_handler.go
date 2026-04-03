package handler

import (
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
func (h *ContestHandler) ListContests(c *gin.Context) {
	p := pagination.GetPagination(c)
	filters := make(map[string]interface{})

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}

	contests, total, err := h.contestService.ListContests(c.Request.Context(), p.Offset(), p.PageSize, filters)
	if err != nil {
		response.InternalError(c, "failed to get contests")
		return
	}

	response.PaginatedResponse(c, contests, total, p.Page, p.PageSize)
}

// GetContest returns a contest by ID.
func (h *ContestHandler) GetContest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	contest, err := h.contestService.GetContest(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "contest not found")
		return
	}

	response.Success(c, contest)
}

// Signup signs up the current user for a contest.
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

	if err := h.contestService.Signup(c.Request.Context(), uint(contestID), userID.(uint)); err != nil {
		response.InternalError(c, "failed to sign up for contest")
		return
	}

	response.Success(c, nil)
}

// Withdraw withdraws the current user from a contest.
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

	if err := h.contestService.Withdraw(c.Request.Context(), uint(contestID), userID.(uint)); err != nil {
		response.InternalError(c, "failed to withdraw from contest")
		return
	}

	response.Success(c, nil)
}

// GetContestProblems returns the problems for a contest.
func (h *ContestHandler) GetContestProblems(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	problems, err := h.contestService.GetContestProblems(c.Request.Context(), uint(contestID))
	if err != nil {
		response.InternalError(c, "failed to get contest problems")
		return
	}

	response.Success(c, problems)
}

// GetContestRanking returns the ranking for a contest.
func (h *ContestHandler) GetContestRanking(c *gin.Context) {
	idStr := c.Param("id")
	contestID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid contest id")
		return
	}

	ranking, err := h.contestService.GetContestRanking(c.Request.Context(), uint(contestID))
	if err != nil {
		response.InternalError(c, "failed to get contest ranking")
		return
	}

	response.Success(c, ranking)
}
