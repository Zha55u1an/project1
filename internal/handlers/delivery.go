package handlers

import (
	"go_project/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeliveryAddressRepository struct {
	db *gorm.DB
}

// Конструктор для создания нового репозитория
func NewDeliveryAddressRepository(db *gorm.DB) *DeliveryAddressRepository {
	return &DeliveryAddressRepository{db: db}
}

// Метод для сохранения адреса доставки
func (r *DeliveryAddressRepository) Save(deliveryAddress *models.DeliveryAddress) error {
	return r.db.Create(deliveryAddress).Error
}


// Обработчик для сохранения адреса доставки
func SaveDeliveryAddressHandler(repo *DeliveryAddressRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Starting SaveDeliveryAddressHandler")

		// Структура запроса
		var request struct {
			AddressLine string  `json:"address_line" binding:"required"` // Адресная строка
			City        string  `json:"city" binding:"required"`        // Город
			Latitude    float64 `json:"latitude" binding:"required"`    // Широта
			Longitude   float64 `json:"longitude" binding:"required"`   // Долгота
		}

		// Проверяем JSON запрос
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		log.Printf("Received request: %+v", request)

		// Извлекаем ID пользователя из контекста
		userID, exists := c.Get("userID")
		if !exists {
			log.Println("User not authenticated")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		log.Printf("User ID from context: %v", userID)

		// Преобразуем userID в тип uint
		userIDUint, ok := userID.(uint)
		if !ok {
			log.Printf("Invalid user ID: %v", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}

		// Проверяем, существует ли адрес для пользователя
		var existingAddress models.DeliveryAddress
		result := repo.db.Where("user_id = ?", userIDUint).First(&existingAddress)

		// Создаём или обновляем адрес
		if result.RowsAffected > 0 {
			// Адрес уже существует, обновляем его
			existingAddress.AddressLine = request.AddressLine
			existingAddress.City = request.City
			existingAddress.Latitude = request.Latitude
			existingAddress.Longitude = request.Longitude

			if err := repo.db.Save(&existingAddress).Error; err != nil {
				log.Printf("Error updating delivery address: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update delivery address"})
				return
			}
			log.Println("Delivery address updated successfully")
			c.JSON(http.StatusOK, gin.H{"message": "Delivery address updated successfully"})
		} else {
			// Адрес не существует, создаём новый
			deliveryAddress := models.DeliveryAddress{
				UserID:      userIDUint,
				AddressLine: request.AddressLine,
				City:        request.City,
				Latitude:    request.Latitude,
				Longitude:   request.Longitude,
			}
			log.Printf("Creating new delivery address: %+v", deliveryAddress)

			if err := repo.db.Create(&deliveryAddress).Error; err != nil {
				log.Printf("Error creating delivery address: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save delivery address"})
				return
			}

			log.Println("Delivery address created successfully")
			c.JSON(http.StatusOK, gin.H{"message": "Delivery address created successfully"})
		}
	}
}


func (repo *DeliveryAddressRepository) GetByUserID(userID uint) (*models.DeliveryAddress, error) {
	var address models.DeliveryAddress
	err := repo.db.Where("user_id = ?", userID).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func GetDeliveryAddressHandler(repo *DeliveryAddressRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Starting GetDeliveryAddressHandler")

		// Извлекаем userID из контекста (установлен middleware)
		userID, exists := c.Get("userID")
		if !exists {
			log.Println("User not authenticated")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Преобразуем userID в тип uint
		userIDUint, ok := userID.(uint)
		if !ok {
			log.Printf("Invalid user ID: %v", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}

		// Извлекаем адрес через репозиторий
		address, err := repo.GetByUserID(userIDUint)
		if err != nil {
			log.Printf("Delivery address not found: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Delivery address not found"})
			return
		}

		// Возвращаем адрес в формате JSON
		c.JSON(http.StatusOK, gin.H{
			"address_line": address.AddressLine,
			"city":         address.City,
		})
	}
}