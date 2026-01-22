package models

import (
	"time"
)

type Transaction struct {
	ID             string `gorm:"type:char(36);primaryKey"`
	Amount         int64
	Type           string // normal | reversal | adjustment
	RefTransaction *string
	CreatedBy      uint
	CreatedAt      time.Time
}
