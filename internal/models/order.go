package models

import (

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id" gorm:"index"` // Внешний ключ к таблице User
	Status     string      `json:"status"`
	OrderItems []OrderItem `gorm:"constraint:OnDelete:CASCADE;"`
}

type OrderItem struct {
	gorm.Model
	OrderID  uint    `json:"order_id"`
	ItemID   uint    `json:"item_id"` 
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"` // Цена товара на момент создания заказа
	Item Item `gorm:"foreignKey:ItemID" json:"item"` // Связь с товаром
}


type PickupPoint struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Review struct {
    gorm.Model
	UserID  uint   `json:"user_id" gorm:"not null"`
    ItemID  uint   `json:"item_id"`
    Rating  int    `json:"rating"`
    Comment string `json:"comment"`
	FilledStars    int `json:"filledStars"`
    RemainingStars int `json:"remainingStars"`
	Item Item `gorm:"foreignKey:ItemID" json:"item"` // Связь с товаром
	User    User   `json:"user" gorm:"foreignKey:UserID"`
}
