package models

import (
    "gorm.io/gorm"
)

type Category struct {
    gorm.Model
    Name string 		`json:"name"`
	Products *[]Item `gorm:"foreignKey:CategoryID"`
}

type Item struct {
    gorm.Model
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CategoryID  uint      `json:"category_id"` 
    Category    *Category `gorm:"foreignKey:CategoryID" json:"category"`
    Price       float64   `json:"price"`
    IsAvailable bool      `json:"is_available"`
}



