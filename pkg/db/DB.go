package db

import (
	"fmt"
	"go_project/internal/models"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", Host, Username, Password, Dbname, strconv.Itoa(Port), SSLMode)

  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

  if err != nil {
    panic(err)
  }

  if err := db.AutoMigrate(&models.User{}, &models.Item{}, &models.Order{}, &models.OrderItem{}); err != nil {
    panic(err)
  }

  DB = db
}
