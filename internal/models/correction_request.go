package models

import "time"

type CorrectionRequest struct {
	ID            string    `gorm:"type:char(36);primaryKey"`
	TransactionID string    `gorm:"type:char(36);not null"`
	Reason        string    `gorm:"type:text;not null"`
	Status        string    `gorm:"type:varchar(20);not null"` // pending, approved, rejected
	RequestedBy   uint      `gorm:"not null"`
	CreatedAt     time.Time
}
