package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
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

// createTeamRequest is the request body for creating a team.
type createTeamRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=64"`
	Slogan string `json:"slogan" binding:"omitempty,max=256"`
}

// updateTeamRequest is the request body for updating a team.
type updateTeamRequest struct {
	Name   string `json:"name" binding:"omitempty,min=1,max=64"`
	Slogan string `json:"slogan" binding:"omitempty,max=256"`
	Avatar string `json:"avatar" binding:"omitempty"`
}

// inviteMemberRequest is the request body for inviting a member.
type inviteMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// CreateTeam creates a new team.
// POST /api/v1/teams
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	var req createTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	team := &model.Team{
		Name:   req.Name,
		Slogan: req.Slogan,
	}

	if err := h.teamService.CreateTeam(c.Request.Context(), team, userID.(uint)); err != nil {
		switch err {
		case service.ErrTeamNameExists:
			response.Error(c, http.StatusConflict, "team name already exists")
		case service.ErrAlreadyInTeam:
			response.BadRequest(c, "you are already in a team")
		case service.ErrUserNotFound:
			response.NotFound(c, "user not found")
		default:
			response.InternalError(c, "failed to create team")
		}
		return
	}

	response.Success(c, team)
}

// GetTeam returns a team by ID.
// GET /api/v1/teams/:id
func (h *TeamHandler) GetTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid team id")
		return
	}

	team, err := h.teamService.GetTeam(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrTeamNotFound) {
			response.NotFound(c, "team not found")
			return
		}
		response.InternalError(c, "failed to get team")
		return
	}

	// Load members for the team response
	members, err := h.teamService.GetTeamMembers(c.Request.Context(), uint(id))
	if err == nil {
		team.Members = members
	}

	response.Success(c, team)
}

// UpdateTeam updates a team.
// PUT /api/v1/teams/:id
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	idStr := c.Param("id")
	teamID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid team id")
		return
	}

	var req updateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	team, err := h.teamService.UpdateTeam(
		c.Request.Context(),
		uint(teamID),
		userID.(uint),
		req.Name,
		req.Slogan,
		req.Avatar,
	)
	if err != nil {
		switch err {
		case service.ErrTeamNotFound:
			response.NotFound(c, "team not found")
		case service.ErrNotTeamLeader:
			response.Forbidden(c, "only the team leader can update the team")
		case service.ErrTeamNameExists:
			response.Error(c, http.StatusConflict, "team name already exists")
		default:
			response.InternalError(c, "failed to update team")
		}
		return
	}

	response.Success(c, team)
}

// DeleteTeam deletes a team.
// DELETE /api/v1/teams/:id
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	idStr := c.Param("id")
	teamID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid team id")
		return
	}

	if err := h.teamService.DeleteTeam(c.Request.Context(), uint(teamID), userID.(uint)); err != nil {
		switch err {
		case service.ErrTeamNotFound:
			response.NotFound(c, "team not found")
		case service.ErrNotTeamLeader:
			response.Forbidden(c, "only the team leader can delete the team")
		default:
			response.InternalError(c, "failed to delete team")
		}
		return
	}

	response.Success(c, nil)
}

// ListTeams returns a paginated list of teams.
// GET /api/v1/teams
func (h *TeamHandler) ListTeams(c *gin.Context) {
	p := pagination.GetPagination(c)
	search := c.Query("search")

	var teams []model.Team
	var total int64
	var err error

	if search != "" {
		teams, total, err = h.teamService.ListTeamsWithSearch(c.Request.Context(), p.Offset(), p.PageSize, search)
	} else {
		teams, total, err = h.teamService.ListTeams(c.Request.Context(), p.Offset(), p.PageSize)
	}

	if err != nil {
		response.InternalError(c, "failed to get teams")
		return
	}

	response.PaginatedResponse(c, teams, total, p.Page, p.PageSize)
}

// InviteUser invites a user to a team.
// POST /api/v1/teams/:id/invite
func (h *TeamHandler) InviteUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	idStr := c.Param("id")
	teamID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid team id")
		return
	}

	var req inviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request parameters: "+err.Error())
		return
	}

	if err := h.teamService.InviteUser(
		c.Request.Context(),
		uint(teamID),
		userID.(uint),
		req.UserID,
	); err != nil {
		switch err {
		case service.ErrTeamNotFound:
			response.NotFound(c, "team not found")
		case service.ErrNotTeamLeader:
			response.Forbidden(c, "only the team leader can invite members")
		case service.ErrCannotInviteSelf:
			response.BadRequest(c, "cannot invite yourself")
		case service.ErrUserNotFound:
			response.NotFound(c, "user not found")
		case service.ErrInviteeInTeam:
			response.BadRequest(c, "the user is already in a team")
		case service.ErrTeamFull:
			response.BadRequest(c, "team has reached the maximum member limit")
		default:
			response.InternalError(c, "failed to invite user")
		}
		return
	}

	response.Success(c, nil)
}

// AcceptInvitation accepts a team invitation.
// POST /api/v1/teams/invitations/:invitationId/accept
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
		switch err {
		case service.ErrInvitationNotFound:
			response.NotFound(c, "invitation not found")
		case service.ErrUnauthorized:
			response.Forbidden(c, "you are not the invitee of this invitation")
		case service.ErrInvitationExpired:
			response.BadRequest(c, "invitation has already been processed")
		case service.ErrAlreadyInTeam:
			response.BadRequest(c, "you are already in a team")
		case service.ErrTeamNotFound:
			response.NotFound(c, "team not found")
		case service.ErrTeamFull:
			response.BadRequest(c, "team has reached the maximum member limit")
		default:
			response.InternalError(c, "failed to accept invitation")
		}
		return
	}

	response.Success(c, nil)
}

// RejectInvitation rejects a team invitation.
// POST /api/v1/teams/invitations/:invitationId/reject
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
		switch err {
		case service.ErrInvitationNotFound:
			response.NotFound(c, "invitation not found")
		case service.ErrUnauthorized:
			response.Forbidden(c, "you are not the invitee of this invitation")
		case service.ErrInvitationExpired:
			response.BadRequest(c, "invitation has already been processed")
		default:
			response.InternalError(c, "failed to reject invitation")
		}
		return
	}

	response.Success(c, nil)
}

// LeaveTeam allows a user to leave a team.
// POST /api/v1/teams/:id/leave
func (h *TeamHandler) LeaveTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	idStr := c.Param("id")
	teamID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid team id")
		return
	}

	if err := h.teamService.LeaveTeam(c.Request.Context(), uint(teamID), userID.(uint)); err != nil {
		switch err {
		case service.ErrTeamNotFound:
			response.NotFound(c, "team not found")
		case service.ErrCannotLeaveAsLeader:
			response.BadRequest(c, "leader cannot leave the team; transfer leadership or disband the team")
		case service.ErrNotTeamMember:
			response.BadRequest(c, "you are not a member of this team")
		default:
			response.InternalError(c, "failed to leave team")
		}
		return
	}

	response.Success(c, nil)
}

// RemoveMember removes a member from a team (leader only).
// DELETE /api/v1/teams/:id/members/:userId
func (h *TeamHandler) RemoveMember(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	idStr := c.Param("id")
	teamID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid team id")
		return
	}

	userIDStr := c.Param("userId")
	memberUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.teamService.RemoveMember(
		c.Request.Context(),
		uint(teamID),
		userID.(uint),
		uint(memberUserID),
	); err != nil {
		switch err {
		case service.ErrTeamNotFound:
			response.NotFound(c, "team not found")
		case service.ErrNotTeamLeader:
			response.Forbidden(c, "only the team leader can remove members")
		case service.ErrNotTeamMember:
			response.BadRequest(c, "the user is not a member of this team")
		default:
			response.InternalError(c, "failed to remove member")
		}
		return
	}

	response.Success(c, nil)
}

// GetInvitations returns pending invitations for the current user.
// GET /api/v1/teams/invitations
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
