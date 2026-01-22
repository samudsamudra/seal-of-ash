package models

import "time"

type ForensicAsh struct {
	ID         uint      `gorm:"primaryKey"`
	EntityType string
	EntityID   uint
	Action     string
	Snapshot   []byte
	ActorID    uint
	IP         string
	UserAgent  string
	Hash       string
	PrevHash   string
	CreatedAt  time.Time
}
