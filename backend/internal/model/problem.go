package model

import "time"

type Problem struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Title           string    `gorm:"type:varchar(256);not null" json:"title"`
	Type            string    `gorm:"type:varchar(20);default:programming;not null" json:"type"` // programming, algorithm, ctf
	CTFCategory     string    `gorm:"type:varchar(20)" json:"ctf_category"` // reverse, pwn, web, crypto, forensics, misc, recon, vuln_reproduce
	Difficulty      string    `gorm:"type:varchar(20);default:easy" json:"difficulty"` // easy, medium, hard
	TimeLimitMs     int       `gorm:"default:2000" json:"time_limit_ms"`
	MemoryLimitKb   int       `gorm:"default:262144" json:"memory_limit_kb"`
	Description     string    `gorm:"type:text" json:"description"`
	InputFormat     string    `gorm:"type:text" json:"input_format"`
	OutputFormat    string    `gorm:"type:text" json:"output_format"`
	Hints           string    `gorm:"type:text" json:"hints"`
	Source          string    `gorm:"type:varchar(20);default:original" json:"source"` // original, ai, imported
	SourceDetail    string    `gorm:"type:varchar(512)" json:"source_detail"`
	SubmitCount     int       `gorm:"default:0" json:"submit_count"`
	AcceptedCount   int       `gorm:"default:0" json:"accepted_count"`
	Status          string    `gorm:"type:varchar(20);default:public;not null" json:"status"` // public, private, disabled
	CTFFlag         string    `gorm:"type:varchar(512)" json:"ctf_flag"`
	Attachments     string    `gorm:"type:json" json:"attachments"`
	SpecialJudgePath string   `gorm:"type:varchar(512)" json:"special_judge_path"`
	CreatedBy       uint      `gorm:"index" json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Tags      []Tag        `gorm:"many2many:problem_tags;" json:"tags,omitempty"`
	TestData  []TestData   `gorm:"foreignKey:ProblemID" json:"test_data,omitempty"`
	Samples   []Sample     `gorm:"foreignKey:ProblemID" json:"samples,omitempty"`
	Creator   *User        `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

type TestData struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	ProblemID   uint    `gorm:"not null;index" json:"problem_id"`
	InputFile   string  `gorm:"type:varchar(512)" json:"input_file"`
	OutputFile  string  `gorm:"type:varchar(512)" json:"output_file"`
	ScoreWeight float64 `gorm:"default:1.0" json:"score_weight"`
	IsSample    bool    `gorm:"default:false" json:"is_sample"`
	Generation  string  `gorm:"type:varchar(20);default:manual" json:"generation"` // manual, ai
	CreatedAt   time.Time `json:"created_at"`

	Problem Problem `gorm:"foreignKey:ProblemID" json:"problem,omitempty"`
}

type Sample struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	ProblemID   uint `gorm:"not null;index" json:"problem_id"`
	TestDataID  uint `gorm:"not null" json:"test_data_id"`
	DisplayOrder int  `gorm:"default:0" json:"display_order"`

	Problem  Problem  `gorm:"foreignKey:ProblemID" json:"problem,omitempty"`
	TestData TestData `gorm:"foreignKey:TestDataID" json:"test_data,omitempty"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`
}

type ProblemTag struct {
	ProblemID uint `gorm:"not null;primaryKey" json:"problem_id"`
	TagID     uint `gorm:"not null;primaryKey" json:"tag_id"`

	Problem Problem `gorm:"foreignKey:ProblemID" json:"problem,omitempty"`
	Tag     Tag     `gorm:"foreignKey:TagID" json:"tag,omitempty"`
}
