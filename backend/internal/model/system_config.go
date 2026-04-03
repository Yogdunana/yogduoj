package model

import "time"

type SystemConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ConfigKey   string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"config_key"`
	ConfigValue string    `gorm:"type:text;not null" json:"config_value"`
	Description string    `gorm:"type:varchar(512)" json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LoginAttempt tracks login attempts for security purposes.
type LoginAttempt struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	IPAddress   string    `gorm:"type:varchar(45);not null" json:"ip_address"`
	AttemptTime time.Time `json:"attempt_time"`
	Success     bool      `gorm:"not null" json:"success"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
