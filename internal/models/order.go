package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	UserID     int         `json:"user_id"`
	CreatedAt  time.Time   `json:"created_at"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"-"`
}

type OrderItem struct {
	OrderID  uint `json:"order_id"`
	ItemID   uint `json:"item_id"`
	Quantity int  `json:"quantity"`
}
