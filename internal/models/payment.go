package models

import "time"

type Payment struct {
	ID        uint      `gorm:"primaryKey"`
	OrderID   string    `gorm:"uniqueIndex"`
	Amount    int64     `json:"amount"`
	PaidAt    time.Time `json:"paid_at"`
	CreatedAt time.Time `json:"created_at"`
}
