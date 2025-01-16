package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string  `json:"name"`
	Products *[]Item `gorm:"foreignKey:CategoryID"`
}

type Item struct {
	gorm.Model
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CategoryID  uint        `json:"category_id"`
	Category    *Category   `gorm:"foreignKey:CategoryID" json:"category"`
	Price       float64     `json:"price"`
	IsAvailable bool        `json:"is_available"`
	SellerID    uint        `json:"seller_id"`       // ID продавца, связанного с товаром
    ImagePath   string    `json:"image_path"` // Только одно изображение
	AverageRating float64  `json:"average_rating"`
	ReviewCount int64  		`json:"review_count"`
	Images        []ItemImage `gorm:"foreignKey:ItemID"`
	OrderItems  []OrderItem `gorm:"constraint:OnDelete:CASCADE;"`
}

type ItemImage struct {
	gorm.Model
	ItemID    uint   `json:"item_id" gorm:"not null"`             // Внешний ключ для связи с товаром
	ImagePath string `json:"image_path" gorm:"size:255;not null"` // Ссылка на изображение
}

type RecentlyViewedItem struct {
	gorm.Model
	UserID    uint      `json:"user_id" gorm:"not null"`    // ID пользователя
	ItemID    uint      `json:"item_id" gorm:"not null"`    // ID товара
	ViewTime  time.Time `json:"view_time" gorm:"not null"`  // Время просмотра товара
	Item      *Item     `gorm:"foreignKey:ItemID"`          // Связь с таблицей товаров
}