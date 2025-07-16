package models

import "time"

type WhatsappLog struct {
	ID           uint `gorm:"primaryKey"`
	Phone        string
	RestaurantID int
	TableNumber  string
	WaiterName   string
	OrderGUID    string
	Amount       int
	PayURL       string
	CreatedAt    time.Time
}
