package handler

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/Yogdunana/yogduoj/backend/internal/pkg/pagination"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// submissionHub maintains a map of active WebSocket connections for judge status updates.
type submissionHub struct {
	mu      sync.RWMutex
	clients map[uint][]*websocket.Conn // submissionID -> connections
}

var hub = &submissionHub{
	clients: make(map[uint][]*websocket.Conn),
}

// Subscribe adds a WebSocket connection for a submission.
func (h *submissionHub) Subscribe(submissionID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[submissionID] = append(h.clients[submissionID], conn)
}

// Unsubscribe removes a WebSocket connection.
func (h *submissionHub) Unsubscribe(submissionID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns := h.clients[submissionID]
	for i, c := range conns {
		if c == conn {
			h.clients[submissionID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(h.clients[submissionID]) == 0 {
		delete(h.clients, submissionID)
	}
}

// Notify sends a judge status update to all subscribers of a submission.
func (h *submissionHub) Notify(submissionID uint, data interface{}) {
	h.mu.RLock()
	conns := h.clients[submissionID]
	h.mu.RUnlock()

	for _, conn := range conns {
		if err := conn.WriteJSON(data); err != nil {
			conn.Close()
			h.Unsubscribe(submissionID, conn)
		}
	}
}

// NotifyJudgeUpdate is a package-level function that allows the judge service
// to push real-time judge status updates to WebSocket clients.
func NotifyJudgeUpdate(submissionID uint, data interface{}) {
	hub.Notify(submissionID, data)
}

// init registers the WebSocket notification callback with the judge service.
func init() {
	service.SetGlobalNotifyFunc(NotifyJudgeUpdate)
}

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
// POST /api/v1/submissions
// JSON body: { "problem_id": 1, "language": "cpp", "code": "...", "ctf_answer": "" }
func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	var req service.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	ipAddress := c.ClientIP()

	submission, err := h.submissionService.CreateSubmission(
		c.Request.Context(),
		userID.(uint),
		req,
		ipAddress,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProblemNotFound):
			response.NotFound(c, "problem not found")
		case errors.Is(err, service.ErrProblemNotPublic):
			response.Forbidden(c, "problem is not available")
		case errors.Is(err, service.ErrInvalidLanguage):
			response.BadRequest(c, "invalid programming language")
		case errors.Is(err, service.ErrCodeEmpty):
			response.BadRequest(c, "code cannot be empty")
		case errors.Is(err, service.ErrCTFAnswerEmpty):
			response.BadRequest(c, "ctf answer cannot be empty")
		default:
			response.InternalError(c, "failed to create submission")
		}
		return
	}

	response.Success(c, submission)
}

// ListSubmissions returns a paginated list of submissions.
// GET /api/v1/submissions?page=1&page_size=20&problem_id=1&result=AC&language=cpp
func (h *SubmissionHandler) ListSubmissions(c *gin.Context) {
	p := pagination.GetPagination(c)

	filter := service.SubmissionFilter{}

	if problemID := c.Query("problem_id"); problemID != "" {
		if id, err := strconv.ParseUint(problemID, 10, 32); err == nil {
			filter.ProblemID = uint(id)
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

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	isAdmin := false
	if role, exists := c.Get("user_role"); exists {
		isAdmin = role.(string) == "admin"
	}

	submissions, total, err := h.submissionService.ListSubmissions(
		c.Request.Context(),
		filter,
		p.Offset(),
		p.PageSize,
		userID.(uint),
		isAdmin,
	)
	if err != nil {
		response.InternalError(c, "failed to get submissions")
		return
	}

	response.PaginatedResponse(c, submissions, total, p.Page, p.PageSize)
}

// GetSubmission returns a submission by ID.
// GET /api/v1/submissions/:id
func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid submission id")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	isAdmin := false
	if role, exists := c.Get("user_role"); exists {
		isAdmin = role.(string) == "admin"
	}

	submission, err := h.submissionService.GetSubmission(
		c.Request.Context(),
		uint(id),
		userID.(uint),
		isAdmin,
	)
	if err != nil {
		if errors.Is(err, service.ErrSubmissionNotFound) {
			response.NotFound(c, "submission not found")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Forbidden(c, "you can only view your own submissions")
			return
		}
		response.InternalError(c, "failed to get submission")
		return
	}

	response.Success(c, submission)
}

// GetSubmissionCode returns the source code of a submission.
// GET /api/v1/submissions/:id/code
func (h *SubmissionHandler) GetSubmissionCode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid submission id")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	isAdmin := false
	if role, exists := c.Get("user_role"); exists {
		isAdmin = role.(string) == "admin"
	}

	code, err := h.submissionService.GetSubmissionCode(
		c.Request.Context(),
		uint(id),
		userID.(uint),
		isAdmin,
	)
	if err != nil {
		if errors.Is(err, service.ErrSubmissionNotFound) {
			response.NotFound(c, "submission not found")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Forbidden(c, "you can only view your own submissions")
			return
		}
		response.InternalError(c, "failed to get submission code")
		return
	}

	response.Success(c, gin.H{
		"code": code,
	})
}

// JudgeWebSocket handles WebSocket connection for real-time judge status.
// GET /api/v1/submissions/:id/judge
func (h *SubmissionHandler) JudgeWebSocket(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid submission id")
		return
	}

	submissionID := uint(id)

	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Register this connection
	hub.Subscribe(submissionID, conn)
	defer hub.Unsubscribe(submissionID, conn)

	// Send current status immediately
	submission, err := h.submissionService.GetSubmission(c.Request.Context(), submissionID, 0, true)
	if err == nil {
		_ = conn.WriteJSON(gin.H{
			"submission_id": submissionID,
			"judge_result":  submission.JudgeResult,
			"judge_score":   submission.JudgeScore,
			"time_used_ms":  submission.TimeUsedMs,
			"memory_used_kb": submission.MemoryUsedKb,
			"error_message": submission.ErrorMessage,
		})
	}

	// Keep connection alive, waiting for close or messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// Client messages are ignored; this is a one-way notification channel
	}
}
