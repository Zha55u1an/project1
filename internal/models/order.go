package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id" gorm:"index"` // Внешний ключ к таблице User
	CreatedAt  time.Time   `json:"created_at"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	OrderID  uint `json:"order_id"`
	ItemID   uint `gorm:"foreignKey:item_id" json:"item_id"`
	Quantity int  `json:"quantity"`
}
