package model

import "time"

type ImportRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	SourcePlatform string   `gorm:"type:varchar(50);not null" json:"source_platform"`
	FilePath      string    `gorm:"type:varchar(512)" json:"file_path"`
	ImportedBy    uint      `gorm:"index" json:"imported_by"`
	ProblemCount  int       `gorm:"default:0" json:"problem_count"`
	SuccessCount  int       `gorm:"default:0" json:"success_count"`
	FailCount     int       `gorm:"default:0" json:"fail_count"`
	Status        string    `gorm:"type:varchar(20);default:pending" json:"status"` // pending, processing, completed, failed
	ErrorMessage  string    `gorm:"type:text" json:"error_message"`
	CreatedAt     time.Time `json:"created_at"`

	Importer *User `gorm:"foreignKey:ImportedBy" json:"importer,omitempty"`
}
