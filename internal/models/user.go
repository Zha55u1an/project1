package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string   `json:"username" gorm:"not null"` // Уникальное имя пользователя
    Email    string   `json:"email" gorm:"unique;"`   // Уникальный email
    Password string   `json:"password" gorm:"not null"`       // Хэшированный пароль
    Role     string   `json:"role" gorm:"size:50;not null"`   // Роль пользователя (например, user, admin)
    IsVerified bool   `json:"is_verified" gorm:"default:false"` // Флаг подтверждения email
    UserInfo  UserInfo `json:"user_info" gorm:"foreignKey:UserID"` // Связь с таблицей UserInfo
    Reviews     []Review  `json:"reviews" gorm:"foreignKey:UserID"`
}

// UserInfo - дополнительная информация о пользователе
type UserInfo struct {
    gorm.Model
    UserID    uint   `json:"user_id" gorm:"not null;unique"` // Внешний ключ
    FirstName string `json:"first_name" gorm:"size:100;not null"`
    LastName  string `json:"last_name" gorm:"size:100;not null"`
    Email     string `json:"email" gorm:"unique;not null"`
    Phone     string `json:"phone" gorm:"size:15;not null"`
}

// DeliveryAddress - адрес доставки
type DeliveryAddress struct {
    gorm.Model
    UserID      uint    `json:"user_id" gorm:"not null"`                  // Связь с пользователем
    AddressLine string  `json:"address_line" gorm:"size:255;not null"`   // Адресная строка
    City        string  `json:"city" gorm:"size:100;not null"`           // Город
    Latitude    float64 `json:"latitude" gorm:"not null"`                // Широта
    Longitude   float64 `json:"longitude" gorm:"not null"`               // Долгота
}


type Verification struct {
    gorm.Model                      // Встроенная модель GORM (содержит ID, CreatedAt, UpdatedAt, DeletedAt)
    Email            string         `json:"email" gorm:"unique;not null"`        // Уникальный email
    VerificationCode string         `json:"verification_code" gorm:"not null"`  // Код подтверждения
    Password         string         `json:"password" gorm:"not null"`           // Временное хранилище хешированного пароля
    Role             string         `json:"role" gorm:"not null"`               // Роль пользователя
    ExpiresAt        time.Time      `json:"expires_at" gorm:"not null"`         // Время истечения кода
}
