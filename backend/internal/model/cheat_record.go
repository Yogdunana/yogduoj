package model

import "time"

type CheatRecord struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"index" json:"user_id"`
	TeamID        uint       `gorm:"index" json:"team_id"`
	ContestID     uint       `gorm:"index" json:"contest_id"`
	SubmissionID  uint       `gorm:"index" json:"submission_id"`
	CheatType     string     `gorm:"type:varchar(50);not null" json:"cheat_type"`
	Detail        string     `gorm:"type:text" json:"detail"`
	SimilarityScore float64  `gorm:"default:0" json:"similarity_score"`
	ReviewStatus  string     `gorm:"type:varchar(20);default:pending" json:"review_status"` // pending, confirmed, dismissed
	ReviewedBy    uint       `json:"reviewed_by"`
	ReviewedAt    *time.Time `json:"reviewed_at"`
	Penalty       string     `gorm:"type:varchar(256)" json:"penalty"`
	CreatedAt     time.Time  `json:"created_at"`

	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Team       Team       `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Contest    Contest    `gorm:"foreignKey:ContestID" json:"contest,omitempty"`
	Submission Submission `gorm:"foreignKey:SubmissionID" json:"submission,omitempty"`
	Reviewer   *User      `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}
