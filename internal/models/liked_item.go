package models

import "gorm.io/gorm"

type LikedItem struct {
	gorm.Model
    ID        uint     `gorm:"primaryKey"`
    UserID    uint     `gorm:"not null"` // Владелец избранных
    ItemID    uint     `gorm:"not null"` // Идентификатор товара
    Item      Item     `gorm:"foreignKey:ItemID"`
}
