package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *model.Team, leaderID uint) error
	GetTeam(ctx context.Context, id uint) (*model.Team, error)
	UpdateTeam(ctx context.Context, team *model.Team) error
	DeleteTeam(ctx context.Context, id uint) error
	ListTeams(ctx context.Context, offset, limit int) ([]model.Team, int64, error)
	InviteUser(ctx context.Context, teamID, inviterID, inviteeID uint) error
	AcceptInvitation(ctx context.Context, invitationID, userID uint) error
	RejectInvitation(ctx context.Context, invitationID, userID uint) error
	LeaveTeam(ctx context.Context, teamID, userID uint) error
	GetTeamMembers(ctx context.Context, teamID uint) ([]model.TeamMember, error)
	GetUserTeams(ctx context.Context, userID uint) ([]model.Team, error)
	GetPendingInvitations(ctx context.Context, userID uint) ([]model.TeamInvitation, error)
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

func (s *teamService) CreateTeam(ctx context.Context, team *model.Team, leaderID uint) error {
	// TODO: implement team creation with member validation
	return s.teamRepo.Create(ctx, team)
}

func (s *teamService) GetTeam(ctx context.Context, id uint) (*model.Team, error) {
	return s.teamRepo.GetByID(ctx, id)
}

func (s *teamService) UpdateTeam(ctx context.Context, team *model.Team) error {
	return s.teamRepo.Update(ctx, team)
}

func (s *teamService) DeleteTeam(ctx context.Context, id uint) error {
	return s.teamRepo.Delete(ctx, id)
}

func (s *teamService) ListTeams(ctx context.Context, offset, limit int) ([]model.Team, int64, error) {
	return s.teamRepo.List(ctx, offset, limit)
}

func (s *teamService) InviteUser(ctx context.Context, teamID, inviterID, inviteeID uint) error {
	// TODO: implement invitation logic with validation
	invitation := &model.TeamInvitation{
		TeamID:    teamID,
		InviterID: inviterID,
		InviteeID: inviteeID,
		Status:    "pending",
	}
	return s.teamRepo.CreateInvitation(ctx, invitation)
}

func (s *teamService) AcceptInvitation(ctx context.Context, invitationID, userID uint) error {
	invitation, err := s.teamRepo.GetInvitationByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation.InviteeID != userID {
		return ErrUnauthorized
	}
	invitation.Status = "accepted"
	if err := s.teamRepo.UpdateInvitation(ctx, invitation); err != nil {
		return err
	}
	member := &model.TeamMember{
		TeamID: invitation.TeamID,
		UserID: userID,
		Role:   "member",
	}
	return s.teamRepo.AddMember(ctx, member)
}

func (s *teamService) RejectInvitation(ctx context.Context, invitationID, userID uint) error {
	invitation, err := s.teamRepo.GetInvitationByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation.InviteeID != userID {
		return ErrUnauthorized
	}
	invitation.Status = "rejected"
	return s.teamRepo.UpdateInvitation(ctx, invitation)
}

func (s *teamService) LeaveTeam(ctx context.Context, teamID, userID uint) error {
	return s.teamRepo.RemoveMember(ctx, teamID, userID)
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
