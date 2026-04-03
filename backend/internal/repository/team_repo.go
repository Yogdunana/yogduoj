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
	ListWithSearch(ctx context.Context, offset, limit int, search string) ([]model.Team, int64, error)
	AddMember(ctx context.Context, member *model.TeamMember) error
	RemoveMember(ctx context.Context, teamID, userID uint) error
	GetMembers(ctx context.Context, teamID uint) ([]model.TeamMember, error)
	GetMemberRole(ctx context.Context, teamID, userID uint) (string, error)
	UpdateMemberRole(ctx context.Context, teamID, userID uint, role string) error
	GetUserTeams(ctx context.Context, userID uint) ([]model.Team, error)
	GetTeamByUserID(ctx context.Context, userID uint) (*model.Team, error)
	CountTeamMembers(ctx context.Context, teamID uint) (int64, error)
	CreateInvitation(ctx context.Context, invitation *model.TeamInvitation) error
	GetInvitationByID(ctx context.Context, id uint) (*model.TeamInvitation, error)
	GetInvitationsByUser(ctx context.Context, userID uint) ([]model.TeamInvitation, error)
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

// ListWithSearch returns a paginated list of teams with optional search by name.
func (r *teamRepository) ListWithSearch(ctx context.Context, offset, limit int, search string) ([]model.Team, int64, error) {
	var teams []model.Team
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Team{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ?", like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Leader").Offset(offset).Limit(limit).
		Order("created_at DESC").Find(&teams).Error; err != nil {
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

// GetMemberRole returns the role of a specific member in a team.
func (r *teamRepository) GetMemberRole(ctx context.Context, teamID, userID uint) (string, error) {
	var member model.TeamMember
	if err := r.db.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		First(&member).Error; err != nil {
		return "", err
	}
	return member.Role, nil
}

// UpdateMemberRole updates the role of a specific member in a team.
func (r *teamRepository) UpdateMemberRole(ctx context.Context, teamID, userID uint, role string) error {
	return r.db.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Update("role", role).Error
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

// GetTeamByUserID returns the team that a user currently belongs to.
// Returns gorm.ErrRecordNotFound if the user is not in any team.
func (r *teamRepository) GetTeamByUserID(ctx context.Context, userID uint) (*model.Team, error) {
	var team model.Team
	if err := r.db.WithContext(ctx).
		Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ? AND teams.status = ?", userID, "active").
		First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

// CountTeamMembers returns the number of members in a team.
func (r *teamRepository) CountTeamMembers(ctx context.Context, teamID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.TeamMember{}).
		Where("team_id = ?", teamID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *teamRepository) CreateInvitation(ctx context.Context, invitation *model.TeamInvitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *teamRepository) GetInvitationByID(ctx context.Context, id uint) (*model.TeamInvitation, error) {
	var invitation model.TeamInvitation
	if err := r.db.WithContext(ctx).Preload("Team").Preload("Inviter").First(&invitation, id).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

// GetInvitationsByUser returns all invitations (any status) for a user.
func (r *teamRepository) GetInvitationsByUser(ctx context.Context, userID uint) ([]model.TeamInvitation, error) {
	var invitations []model.TeamInvitation
	if err := r.db.WithContext(ctx).
		Where("invitee_id = ?", userID).
		Preload("Team").Preload("Inviter").
		Order("created_at DESC").
		Find(&invitations).Error; err != nil {
		return nil, err
	}
	return invitations, nil
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
