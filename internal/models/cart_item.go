package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
    ID        uint     `gorm:"primaryKey"`
    UserID    uint     `gorm:"not null"` // Владелец корзины
    ItemID    uint     `gorm:"not null"` // Идентификатор товара
    Quantity  int      `gorm:"default:1"`
    Item      Item     `gorm:"foreignKey:ItemID"`
}
