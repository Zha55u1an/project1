package handlers

import (
	"go_project/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Обработчик для отображения профиля пользователя
// Обработчик для отображения профиля пользователя
func ProfileHandler(repo *UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID текущего пользователя из контекста
		userID, exists := c.Get("userID")
		if !exists {
			log.Printf("User ID not found in context\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var user models.User
		var UserInfo models.UserInfo

		// Загружаем данные пользователя по ID
		log.Printf("Fetching user with ID: %v\n", userID)
		result := repo.db.Where("id = ?", userID).First(&user)
		if result.Error != nil {
			log.Printf("Error fetching user: %v\n", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		// Загружаем дополнительную информацию о пользователе (например, имя, фамилию и т.д.)
		log.Printf("Fetching user info for userID: %v\n", userID)
		result = repo.db.Where("user_id = ?", userID).First(&UserInfo)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Если данных о пользователе нет, отправляем пустую структуру или сообщение
				log.Printf("No user info found for userID: %v\n", userID)
				c.HTML(http.StatusOK, "profile.html", gin.H{
					"user":     user,
					"message":  "Пожалуйста, заполните информацию о профиле",
				})
				return
			}
			log.Printf("Error fetching user info: %v\n", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		log.Printf("UserInfo data: %+v", UserInfo)

		// Отправляем данные профиля на страницу
		log.Printf("Successfully fetched user info for userID: %v\n", userID)
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"user":     user,
			"UserInfo": UserInfo,
		})
	}
}



// Обработчик для обновления данных профиля пользователя
func UpdateProfileHandler(repo *UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID текущего пользователя из контекста
		userID, exists := c.Get("userID")
		if !exists {
			// Логируем ошибку
			log.Printf("Error: Unauthorized access attempt\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var userInfo models.UserInfo

		// Привязываем данные из формы
		if err := c.ShouldBindJSON(&userInfo); err != nil {
			// Логируем ошибку
			log.Printf("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Логируем перед обновлением
		log.Printf("Attempting to update user info for userID: %v, New Data: %+v\n", userID, userInfo)

		// Обновляем информацию о пользователе через репозиторий
		err := repo.UpdateUserInfo(userID.(uint), &userInfo)
		if err != nil {
			// Логируем ошибку при обновлении
			log.Printf("Error updating user info for userID: %v, Error: %v\n", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
			return
		}

		// Логируем успешное обновление
		log.Printf("Successfully updated user info for userID: %v\n", userID)
		c.JSON(http.StatusOK, gin.H{"message": "User info updated successfully"})
	}
}



func (repo *UserRepository) UpdateUserInfo(userID uint, userInfo *models.UserInfo) error {
	// Логируем начало обновления
	log.Printf("Updating user info for userID: %v, New Data: %+v\n", userID, userInfo)

	// Сначала пытаемся найти существующую запись для данного пользователя
	var existingUserInfo models.UserInfo
	if err := repo.db.Where("user_id = ?", userID).First(&existingUserInfo).Error; err != nil {
		// Если записи нет, создаем новую
		if err == gorm.ErrRecordNotFound {
			log.Printf("No record found for userID %v, creating new record.\n", userID)

			// Заполняем данные для создания новой записи
			userInfo.UserID = userID

			// Создаем новую запись
			if err := repo.db.Create(userInfo).Error; err != nil {
				log.Printf("Error creating new user info for userID: %v, Error: %v\n", userID, err)
				return err
			}

			log.Printf("Successfully created new user info for userID: %v\n", userID)
			return nil
		} else {
			// Логируем ошибку при поиске записи
			log.Printf("Error searching for user info for userID: %v, Error: %v\n", userID, err)
			return err
		}
	}

	// Если запись найдена, обновляем ее
	log.Printf("Found existing record for userID: %v, updating data.\n", userID)

	// Обновляем информацию о пользователе в базе данных
	result := repo.db.Model(&existingUserInfo).Updates(userInfo)
	if result.Error != nil {
		// Логируем ошибку при обновлении
		log.Printf("Error updating user info in database: %v\n", result.Error)
		return result.Error // Если произошла ошибка, возвращаем её
	}

	// Логируем успешное завершение обновления
	log.Printf("Successfully updated user info for userID: %v\n", userID)
	return nil // Если все прошло успешно, возвращаем nil (без ошибок)
}


