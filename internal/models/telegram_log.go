package models

import "time"

type TelegramLog struct {
	ID           uint `gorm:"primaryKey"`
	TelegramID   int64
	RestaurantID int
	TableNumber  string
	WaiterName   string
	OrderGUID    string
	Amount       int
	PayURL       string
	CreatedAt    time.Time
}
