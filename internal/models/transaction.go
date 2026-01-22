package models

import "time"

type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	Amount          int64
	Type            string    // normal, reversal, adjustment
	RefTransaction  *uint
	CreatedBy       uint
	CreatedAt       time.Time
}
