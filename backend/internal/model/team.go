package model

import "time"

type Team struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`
	Slogan     string    `gorm:"type:varchar(256)" json:"slogan"`
	LeaderID   uint      `gorm:"not null;index" json:"leader_id"`
	Avatar     string    `gorm:"type:varchar(512)" json:"avatar"`
	MaxMembers int       `gorm:"default:3" json:"max_members"`
	Status     string    `gorm:"type:varchar(20);default:active;not null" json:"status"` // active, disbanded
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Leader      User         `gorm:"foreignKey:LeaderID" json:"leader,omitempty"`
	Members     []TeamMember `gorm:"foreignKey:TeamID" json:"members,omitempty"`
}

type TeamMember struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TeamID      uint      `gorm:"not null;uniqueIndex:idx_team_user" json:"team_id"`
	UserID      uint      `gorm:"not null;uniqueIndex:idx_team_user" json:"user_id"`
	Role        string    `gorm:"type:varchar(20);default:member;not null" json:"role"` // leader, member
	Permission  string    `gorm:"type:json" json:"permission"`
	JoinedAt    time.Time `json:"joined_at"`

	Team Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type TeamInvitation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TeamID     uint      `gorm:"not null;index" json:"team_id"`
	InviterID  uint      `gorm:"not null" json:"inviter_id"`
	InviteeID  uint      `gorm:"not null;index" json:"invitee_id"`
	Status     string    `gorm:"type:varchar(20);default:pending;not null" json:"status"` // pending, accepted, rejected
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Team    Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Inviter User `gorm:"foreignKey:InviterID" json:"inviter,omitempty"`
	Invitee User `gorm:"foreignKey:InviteeID" json:"invitee,omitempty"`
}
