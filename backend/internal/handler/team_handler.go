package handler

import (
	"strconv"

	"github.com/Yogdunana/yogduoj/backend/internal/pkg/pagination"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// CreateTeam creates a new team.
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	// TODO: implement team creation
	response.Error(c, 501, "not implemented")
}

// GetTeam returns a team by ID.
func (h *TeamHandler) GetTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid team id")
		return
	}

	team, err := h.teamService.GetTeam(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "team not found")
		return
	}

	response.Success(c, team)
}

// UpdateTeam updates a team.
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// DeleteTeam deletes a team.
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// ListTeams returns a paginated list of teams.
func (h *TeamHandler) ListTeams(c *gin.Context) {
	p := pagination.GetPagination(c)
	teams, total, err := h.teamService.ListTeams(c.Request.Context(), p.Offset(), p.PageSize)
	if err != nil {
		response.InternalError(c, "failed to get teams")
		return
	}

	response.PaginatedResponse(c, teams, total, p.Page, p.PageSize)
}

// InviteUser invites a user to a team.
func (h *TeamHandler) InviteUser(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// AcceptInvitation accepts a team invitation.
func (h *TeamHandler) AcceptInvitation(c *gin.Context) {
	idStr := c.Param("invitationId")
	invitationID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid invitation id")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	if err := h.teamService.AcceptInvitation(c.Request.Context(), uint(invitationID), userID.(uint)); err != nil {
		response.InternalError(c, "failed to accept invitation")
		return
	}

	response.Success(c, nil)
}

// RejectInvitation rejects a team invitation.
func (h *TeamHandler) RejectInvitation(c *gin.Context) {
	idStr := c.Param("invitationId")
	invitationID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid invitation id")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	if err := h.teamService.RejectInvitation(c.Request.Context(), uint(invitationID), userID.(uint)); err != nil {
		response.InternalError(c, "failed to reject invitation")
		return
	}

	response.Success(c, nil)
}

// LeaveTeam allows a user to leave a team.
func (h *TeamHandler) LeaveTeam(c *gin.Context) {
	// TODO: implement
	response.Error(c, 501, "not implemented")
}

// GetInvitations returns pending invitations for the current user.
func (h *TeamHandler) GetInvitations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	invitations, err := h.teamService.GetPendingInvitations(c.Request.Context(), userID.(uint))
	if err != nil {
		response.InternalError(c, "failed to get invitations")
		return
	}

	response.Success(c, invitations)
}
