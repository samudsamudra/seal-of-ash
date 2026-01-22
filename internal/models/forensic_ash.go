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
	Hash       string
	PrevHash   string
	CreatedAt  time.Time
}
