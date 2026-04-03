package model

import "time"

type Contest struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Title             string     `gorm:"type:varchar(256);not null" json:"title"`
	ContestType       string     `gorm:"type:varchar(20);default:individual;not null" json:"contest_type"` // individual, team
	Category          string     `gorm:"type:varchar(20);default:programming" json:"category"` // programming, algorithm, ai_assisted, ctf
	RuleType          string     `gorm:"type:varchar(20);default:acm" json:"rule_type"` // acm, oi, ioi, cf, ctf, awd, isw, diy
	DIYTemplateID     uint       `gorm:"index" json:"diy_template_id"`
	StartTime         time.Time  `json:"start_time"`
	EndTime           time.Time  `json:"end_time"`
	FreezeTime        *time.Time `json:"freeze_time"`
	Description       string     `gorm:"type:text" json:"description"`
	RuleDescription   string     `gorm:"type:text" json:"rule_description"`
	MaxTeamSize       int        `gorm:"default:3" json:"max_team_size"`
	SignupLimit       int        `gorm:"default:0" json:"signup_limit"` // 0 means no limit
	AllowViewOthers   bool       `gorm:"default:true" json:"allow_view_others"`
	ShowRealtimeRank  bool       `gorm:"default:true" json:"show_realtime_rank"`
	EnableAIHint      bool       `gorm:"default:false" json:"enable_ai_hint"`
	DIYRules          string     `gorm:"type:json" json:"diy_rules"`
	Status            string     `gorm:"type:varchar(20);default:pending;index" json:"status"` // pending, running, ended, cancelled
	ParticipantCount  int        `gorm:"default:0" json:"participant_count"`
	CreatedBy         uint       `gorm:"index" json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	Creator     *User             `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Problems    []ContestProblem  `gorm:"foreignKey:ContestID" json:"problems,omitempty"`
	Signups     []ContestSignup   `gorm:"foreignKey:ContestID" json:"signups,omitempty"`
}

type ContestProblem struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	ContestID    uint    `gorm:"not null;index" json:"contest_id"`
	ProblemID    uint    `gorm:"not null" json:"problem_id"`
	DisplayOrder int     `gorm:"default:0" json:"display_order"`
	Score        float64 `gorm:"default:0" json:"score"`
	ProblemLabel string  `gorm:"type:varchar(10)" json:"problem_label"`

	Contest Contest `gorm:"foreignKey:ContestID" json:"contest,omitempty"`
	Problem Problem `gorm:"foreignKey:ProblemID" json:"problem,omitempty"`
}

type ContestSignup struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ContestID  uint      `gorm:"not null;uniqueIndex:idx_contest_user" json:"contest_id"`
	UserID     uint      `gorm:"not null;uniqueIndex:idx_contest_user" json:"user_id"`
	TeamID     uint      `gorm:"index" json:"team_id"`
	SignupTime time.Time `json:"signup_time"`

	Contest Contest `gorm:"foreignKey:ContestID" json:"contest,omitempty"`
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Team    Team    `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

type DIYContestTemplate struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(128);not null" json:"name"`
	ScoringRule  string    `gorm:"type:json" json:"scoring_rule"`
	PenaltyRule  string    `gorm:"type:json" json:"penalty_rule"`
	RankingRule  string    `gorm:"type:json" json:"ranking_rule"`
	Description  string    `gorm:"type:text" json:"description"`
	CreatedBy    uint      `gorm:"index" json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`

	Creator *User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}
