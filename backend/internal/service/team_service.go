package service

import (
	"context"
	"errors"
	"time"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrTeamNameExists     = errors.New("team name already exists")
	ErrTeamNotFound       = errors.New("team not found")
	ErrNotTeamLeader      = errors.New("only the team leader can perform this action")
	ErrAlreadyInTeam      = errors.New("user is already in a team")
	ErrNotTeamMember      = errors.New("user is not a member of this team")
	ErrTeamFull           = errors.New("team has reached maximum member limit")
	ErrCannotLeaveAsLeader = errors.New("leader cannot leave the team; transfer leadership or disband the team")
	ErrInvitationNotFound = errors.New("invitation not found")
	ErrInvitationExpired  = errors.New("invitation has already been processed")
	ErrCannotInviteSelf   = errors.New("cannot invite yourself")
	ErrInviteeInTeam      = errors.New("invitee is already in a team")
	ErrAlreadyInvited     = errors.New("invitation already exists")
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *model.Team, leaderID uint) error
	GetTeam(ctx context.Context, id uint) (*model.Team, error)
	UpdateTeam(ctx context.Context, teamID, leaderID uint, name, slogan, avatar string) (*model.Team, error)
	DeleteTeam(ctx context.Context, teamID, leaderID uint) error
	ListTeams(ctx context.Context, offset, limit int) ([]model.Team, int64, error)
	ListTeamsWithSearch(ctx context.Context, offset, limit int, search string) ([]model.Team, int64, error)
	InviteUser(ctx context.Context, teamID, inviterID, inviteeID uint) error
	AcceptInvitation(ctx context.Context, invitationID, userID uint) error
	RejectInvitation(ctx context.Context, invitationID, userID uint) error
	LeaveTeam(ctx context.Context, teamID, userID uint) error
	RemoveMember(ctx context.Context, teamID, leaderID, memberUserID uint) error
	GetTeamMembers(ctx context.Context, teamID uint) ([]model.TeamMember, error)
	GetUserTeams(ctx context.Context, userID uint) ([]model.Team, error)
	GetPendingInvitations(ctx context.Context, userID uint) ([]model.TeamInvitation, error)
	ValidateTeamMembership(ctx context.Context, userID, teamID uint) (bool, error)
}

type teamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// CreateTeam creates a new team with the given user as leader.
// Validates that the team name is unique and the user is not already in another team.
func (s *teamService) CreateTeam(ctx context.Context, team *model.Team, leaderID uint) error {
	// Check if team name already exists
	existing, err := s.teamRepo.GetByName(ctx, team.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return ErrTeamNameExists
	}

	// Check if the user is already in a team
	_, err = s.teamRepo.GetTeamByUserID(ctx, leaderID)
	if err == nil {
		// User is already in a team (no error means record found)
		return ErrAlreadyInTeam
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Verify the leader user exists
	_, err = s.userRepo.GetByID(ctx, leaderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Set team fields
	team.LeaderID = leaderID
	team.Status = "active"
	if team.MaxMembers <= 0 {
		team.MaxMembers = 3
	}

	// Create the team
	if err := s.teamRepo.Create(ctx, team); err != nil {
		return err
	}

	// Add the leader as a team member with leader role
	member := &model.TeamMember{
		TeamID:   team.ID,
		UserID:   leaderID,
		Role:     "leader",
		JoinedAt: time.Now(),
	}
	return s.teamRepo.AddMember(ctx, member)
}

func (s *teamService) GetTeam(ctx context.Context, id uint) (*model.Team, error) {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}
	return team, nil
}

// UpdateTeam updates team information. Only the leader can perform this action.
func (s *teamService) UpdateTeam(ctx context.Context, teamID, leaderID uint, name, slogan, avatar string) (*model.Team, error) {
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}

	// Verify the caller is the leader
	if team.LeaderID != leaderID {
		return nil, ErrNotTeamLeader
	}

	// Update name if provided
	if name != "" && name != team.Name {
		existing, err := s.teamRepo.GetByName(ctx, name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil && existing.ID != team.ID {
			return nil, ErrTeamNameExists
		}
		team.Name = name
	}

	// Update slogan if provided
	if slogan != "" {
		team.Slogan = slogan
	}

	// Update avatar if provided
	if avatar != "" {
		team.Avatar = avatar
	}

	if err := s.teamRepo.Update(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

// DeleteTeam deletes a team. Only the leader can perform this action.
// All members are removed and invitations are cleaned up.
func (s *teamService) DeleteTeam(ctx context.Context, teamID, leaderID uint) error {
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTeamNotFound
		}
		return err
	}

	if team.LeaderID != leaderID {
		return ErrNotTeamLeader
	}

	// Mark team as disbanded
	team.Status = "disbanded"
	if err := s.teamRepo.Update(ctx, team); err != nil {
		return err
	}

	// Remove all members
	members, err := s.teamRepo.GetMembers(ctx, teamID)
	if err != nil {
		return err
	}
	for _, m := range members {
		if err := s.teamRepo.RemoveMember(ctx, teamID, m.UserID); err != nil {
			return err
		}
	}

	// Delete the team record
	return s.teamRepo.Delete(ctx, teamID)
}

func (s *teamService) ListTeams(ctx context.Context, offset, limit int) ([]model.Team, int64, error) {
	return s.teamRepo.List(ctx, offset, limit)
}

// ListTeamsWithSearch returns a paginated list of teams with optional search by name.
func (s *teamService) ListTeamsWithSearch(ctx context.Context, offset, limit int, search string) ([]model.Team, int64, error) {
	return s.teamRepo.ListWithSearch(ctx, offset, limit, search)
}

// InviteUser invites a user to join a team. Only the leader can invite.
func (s *teamService) InviteUser(ctx context.Context, teamID, inviterID, inviteeID uint) error {
	// Get team and verify inviter is the leader
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTeamNotFound
		}
		return err
	}

	if team.LeaderID != inviterID {
		return ErrNotTeamLeader
	}

	if team.Status != "active" {
		return ErrTeamNotFound
	}

	// Cannot invite yourself
	if inviterID == inviteeID {
		return ErrCannotInviteSelf
	}

	// Check if invitee exists
	_, err = s.userRepo.GetByID(ctx, inviteeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Check if invitee is already in a team
	_, err = s.teamRepo.GetTeamByUserID(ctx, inviteeID)
	if err == nil {
		return ErrInviteeInTeam
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check if team is full
	memberCount, err := s.teamRepo.CountTeamMembers(ctx, teamID)
	if err != nil {
		return err
	}
	if int(memberCount) >= team.MaxMembers {
		return ErrTeamFull
	}

	// Create invitation
	invitation := &model.TeamInvitation{
		TeamID:    teamID,
		InviterID: inviterID,
		InviteeID: inviteeID,
		Status:    "pending",
	}
	return s.teamRepo.CreateInvitation(ctx, invitation)
}

// AcceptInvitation allows a user to accept a pending team invitation.
func (s *teamService) AcceptInvitation(ctx context.Context, invitationID, userID uint) error {
	invitation, err := s.teamRepo.GetInvitationByID(ctx, invitationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvitationNotFound
		}
		return err
	}

	if invitation.InviteeID != userID {
		return ErrUnauthorized
	}

	if invitation.Status != "pending" {
		return ErrInvitationExpired
	}

	// Check if user is already in a team
	_, err = s.teamRepo.GetTeamByUserID(ctx, userID)
	if err == nil {
		return ErrAlreadyInTeam
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check if team is still active and has room
	team, err := s.teamRepo.GetByID(ctx, invitation.TeamID)
	if err != nil {
		return ErrTeamNotFound
	}
	if team.Status != "active" {
		return ErrTeamNotFound
	}

	memberCount, err := s.teamRepo.CountTeamMembers(ctx, invitation.TeamID)
	if err != nil {
		return err
	}
	if int(memberCount) >= team.MaxMembers {
		return ErrTeamFull
	}

	// Update invitation status
	invitation.Status = "accepted"
	if err := s.teamRepo.UpdateInvitation(ctx, invitation); err != nil {
		return err
	}

	// Add user as team member
	member := &model.TeamMember{
		TeamID:   invitation.TeamID,
		UserID:   userID,
		Role:     "member",
		JoinedAt: time.Now(),
	}
	return s.teamRepo.AddMember(ctx, member)
}

// RejectInvitation allows a user to reject a pending team invitation.
func (s *teamService) RejectInvitation(ctx context.Context, invitationID, userID uint) error {
	invitation, err := s.teamRepo.GetInvitationByID(ctx, invitationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvitationNotFound
		}
		return err
	}

	if invitation.InviteeID != userID {
		return ErrUnauthorized
	}

	if invitation.Status != "pending" {
		return ErrInvitationExpired
	}

	invitation.Status = "rejected"
	return s.teamRepo.UpdateInvitation(ctx, invitation)
}

// LeaveTeam allows a member (non-leader) to leave a team.
func (s *teamService) LeaveTeam(ctx context.Context, teamID, userID uint) error {
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTeamNotFound
		}
		return err
	}

	if team.LeaderID == userID {
		return ErrCannotLeaveAsLeader
	}

	// Verify the user is a member
	role, err := s.teamRepo.GetMemberRole(ctx, teamID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotTeamMember
		}
		return err
	}
	_ = role

	return s.teamRepo.RemoveMember(ctx, teamID, userID)
}

// RemoveMember allows the team leader to remove a member from the team.
func (s *teamService) RemoveMember(ctx context.Context, teamID, leaderID, memberUserID uint) error {
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTeamNotFound
		}
		return err
	}

	if team.LeaderID != leaderID {
		return ErrNotTeamLeader
	}

	// Cannot remove the leader
	if leaderID == memberUserID {
		return ErrNotTeamLeader
	}

	// Verify the target is a member
	_, err = s.teamRepo.GetMemberRole(ctx, teamID, memberUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotTeamMember
		}
		return err
	}

	return s.teamRepo.RemoveMember(ctx, teamID, memberUserID)
}

func (s *teamService) GetTeamMembers(ctx context.Context, teamID uint) ([]model.TeamMember, error) {
	return s.teamRepo.GetMembers(ctx, teamID)
}

func (s *teamService) GetUserTeams(ctx context.Context, userID uint) ([]model.Team, error) {
	return s.teamRepo.GetUserTeams(ctx, userID)
}

func (s *teamService) GetPendingInvitations(ctx context.Context, userID uint) ([]model.TeamInvitation, error) {
	return s.teamRepo.GetPendingInvitations(ctx, userID)
}

// ValidateTeamMembership checks whether a user is a member of a given team.
func (s *teamService) ValidateTeamMembership(ctx context.Context, userID, teamID uint) (bool, error) {
	_, err := s.teamRepo.GetMemberRole(ctx, teamID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
