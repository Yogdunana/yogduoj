package model

import "time"

type User struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Username         string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Email            string    `gorm:"type:varchar(128);uniqueIndex" json:"email"`
	Phone            string    `gorm:"type:varchar(20)" json:"phone"`
	PasswordHash     string    `gorm:"type:varchar(255);not null" json:"-"`
	Role             string    `gorm:"type:varchar(20);default:user;not null" json:"role"` // user, admin
	Status           string    `gorm:"type:varchar(20);default:active;not null" json:"status"` // active, disabled
	Avatar           string    `gorm:"type:varchar(512)" json:"avatar"`
	SolvedCount      int       `gorm:"default:0" json:"solved_count"`
	SubmissionCount  int       `gorm:"default:0" json:"submission_count"`
	ContestCount     int       `gorm:"default:0" json:"contest_count"`
	Rating           int       `gorm:"default:0" json:"rating"`
	UsernameModified bool      `gorm:"default:false" json:"username_modified"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
