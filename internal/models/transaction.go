package models

import (
	"time"
)

type Transaction struct {
	ID             string `gorm:"type:char(36);primaryKey"`
	RootID         string `gorm:"type:char(36);index"`
	Amount         int64
	Type           string
	RefTransaction *string
	CreatedBy      uint
	CreatedAt      time.Time
}
