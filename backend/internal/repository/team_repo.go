package repository

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	GetByID(ctx context.Context, id uint) (*model.Team, error)
	GetByName(ctx context.Context, name string) (*model.Team, error)
	Update(ctx context.Context, team *model.Team) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]model.Team, int64, error)
	AddMember(ctx context.Context, member *model.TeamMember) error
	RemoveMember(ctx context.Context, teamID, userID uint) error
	GetMembers(ctx context.Context, teamID uint) ([]model.TeamMember, error)
	GetUserTeams(ctx context.Context, userID uint) ([]model.Team, error)
	CreateInvitation(ctx context.Context, invitation *model.TeamInvitation) error
	GetInvitationByID(ctx context.Context, id uint) (*model.TeamInvitation, error)
	UpdateInvitation(ctx context.Context, invitation *model.TeamInvitation) error
	GetPendingInvitations(ctx context.Context, userID uint) ([]model.TeamInvitation, error)
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, team *model.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *teamRepository) GetByID(ctx context.Context, id uint) (*model.Team, error) {
	var team model.Team
	if err := r.db.WithContext(ctx).Preload("Leader").First(&team, id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) GetByName(ctx context.Context, name string) (*model.Team, error) {
	var team model.Team
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) Update(ctx context.Context, team *model.Team) error {
	return r.db.WithContext(ctx).Save(team).Error
}

func (r *teamRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Team{}, id).Error
}

func (r *teamRepository) List(ctx context.Context, offset, limit int) ([]model.Team, int64, error) {
	var teams []model.Team
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.Team{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

func (r *teamRepository) AddMember(ctx context.Context, member *model.TeamMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *teamRepository) RemoveMember(ctx context.Context, teamID, userID uint) error {
	return r.db.WithContext(ctx).Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&model.TeamMember{}).Error
}

func (r *teamRepository) GetMembers(ctx context.Context, teamID uint) ([]model.TeamMember, error) {
	var members []model.TeamMember
	if err := r.db.WithContext(ctx).Where("team_id = ?", teamID).
		Preload("User").Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *teamRepository) GetUserTeams(ctx context.Context, userID uint) ([]model.Team, error) {
	var teams []model.Team
	if err := r.db.WithContext(ctx).
		Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ?", userID).
		Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *teamRepository) CreateInvitation(ctx context.Context, invitation *model.TeamInvitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *teamRepository) GetInvitationByID(ctx context.Context, id uint) (*model.TeamInvitation, error) {
	var invitation model.TeamInvitation
	if err := r.db.WithContext(ctx).First(&invitation, id).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *teamRepository) UpdateInvitation(ctx context.Context, invitation *model.TeamInvitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}

func (r *teamRepository) GetPendingInvitations(ctx context.Context, userID uint) ([]model.TeamInvitation, error) {
	var invitations []model.TeamInvitation
	if err := r.db.WithContext(ctx).
		Where("invitee_id = ? AND status = ?", userID, "pending").
		Preload("Team").Preload("Inviter").
		Find(&invitations).Error; err != nil {
		return nil, err
	}
	return invitations, nil
}
