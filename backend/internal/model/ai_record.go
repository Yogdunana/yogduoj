package model

import "time"

type AIProblemRecord struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ProblemID    uint      `gorm:"index" json:"problem_id"`
	GeneratedBy  uint      `gorm:"index" json:"generated_by"`
	Params       string    `gorm:"type:json" json:"params"`
	AIResponse   string    `gorm:"type:text" json:"ai_response"`
	ReviewStatus string    `gorm:"type:varchar(20);default:pending" json:"review_status"` // pending, approved, rejected
	CreatedAt    time.Time `json:"created_at"`

	Problem   Problem `gorm:"foreignKey:ProblemID" json:"problem,omitempty"`
	Generator *User   `gorm:"foreignKey:GeneratedBy" json:"generator,omitempty"`
}

type AITestdataRecord struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TestDataID   uint      `gorm:"index" json:"test_data_id"`
	ProblemID    uint      `gorm:"index" json:"problem_id"`
	GeneratedBy  uint      `gorm:"index" json:"generated_by"`
	Params       string    `gorm:"type:json" json:"params"`
	ReviewStatus string    `gorm:"type:varchar(20);default:pending" json:"review_status"` // pending, approved, rejected
	CreatedAt    time.Time `json:"created_at"`

	TestData  TestData `gorm:"foreignKey:TestDataID" json:"test_data,omitempty"`
	Problem   Problem  `gorm:"foreignKey:ProblemID" json:"problem,omitempty"`
	Generator *User    `gorm:"foreignKey:GeneratedBy" json:"generator,omitempty"`
}
