package models

import "time"

type ForensicAsh struct {
	ID         string    `gorm:"type:char(36);primaryKey"`
	EntityType string
	EntityID   string
	Action     string
	Snapshot   []byte    `gorm:"type:longblob"`
	ActorID    uint
	IP         string
	UserAgent  string
	Hash       string    `gorm:"size:64;not null"`
	PrevHash   string    `gorm:"size:64;not null"`
	CreatedAt  time.Time
}
