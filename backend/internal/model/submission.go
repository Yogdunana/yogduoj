package model

import "time"

type Submission struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"not null;index" json:"user_id"`
	TeamID        uint       `gorm:"index" json:"team_id"`
	ProblemID     uint       `gorm:"not null;index" json:"problem_id"`
	ContestID     uint       `gorm:"index" json:"contest_id"`
	Language      string     `gorm:"type:varchar(20);not null" json:"language"` // cpp, c, java, python3
	CodePath      string     `gorm:"type:varchar(512)" json:"code_path"`
	CodeLength    int        `gorm:"default:0" json:"code_length"`
	CTFAnswer     string     `gorm:"type:varchar(512)" json:"ctf_answer"`
	CTFFilePath   string     `gorm:"type:varchar(512)" json:"ctf_file_path"`
	JudgeResult   string     `gorm:"type:varchar(20);default:pending;index" json:"judge_result"` // pending, judging, AC, WA, TLE, RE, CE, MLE, SE, WR
	JudgeScore    float64    `gorm:"default:0" json:"judge_score"`
	TimeUsedMs    int        `gorm:"default:0" json:"time_used_ms"`
	MemoryUsedKb  int        `gorm:"default:0" json:"memory_used_kb"`
	ErrorMessage  string     `gorm:"type:text" json:"error_message"`
	JudgeDetail   string     `gorm:"type:json" json:"judge_detail"`
	IPAddress     string     `gorm:"type:varchar(45)" json:"ip_address"`
	SubmitTime    time.Time  `json:"submit_time"`
	JudgeStart    *time.Time `json:"judge_start"`
	JudgeEnd      *time.Time `json:"judge_end"`

	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Team    Team    `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Problem Problem `gorm:"foreignKey:ProblemID" json:"problem,omitempty"`
	Contest Contest `gorm:"foreignKey:ContestID" json:"contest,omitempty"`
}
