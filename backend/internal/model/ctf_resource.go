package model

import "time"

type CTFResource struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(256);not null" json:"title"`
	ResourceType string   `gorm:"type:varchar(20);not null" json:"resource_type"` // tool, tutorial, knowledge, writeup
	CTFCategory string    `gorm:"type:varchar(20)" json:"ctf_category"`
	Content     string    `gorm:"type:text" json:"content"`
	FilePath    string    `gorm:"type:varchar(512)" json:"file_path"`
	ExternalURL string    `gorm:"type:varchar(512)" json:"external_url"`
	CreatedBy   uint      `gorm:"index" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`

	Creator *User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}
