package handlers

import (
	"errors"
	"go_project/internal/models"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LikedRepository struct {
	db *gorm.DB
}

func NewLikedRepository(db *gorm.DB) *LikedRepository {
	return &LikedRepository{db: db}
}

// Обработчик для добавления товара в избранных
func AddToLikedHandler(repo *LikedRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			ItemID   uint `json:"item_id"`
		}

		// Проверяем JSON запрос
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Проверка, что ItemID корректны
        if request.ItemID == 0 {
            log.Println("Invalid item_id in request:", request.ItemID)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_id"})
            return
        }

		// Извлекаем ID пользователя из контекста (должен быть установлен аутентификационным middleware)
		userID, exists := c.Get("userID")
		if !exists {
			log.Println("User not authenticated")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Преобразуем userID в тип uint
		userIDUint, ok := userID.(uint)
		if !ok {
			log.Println("Invalid user ID:", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}

		// Добавляем товар в корзину
		action, err := repo.AddToLiked(userIDUint, request.ItemID)
		if err != nil {
			log.Println("Unable to add item to cart:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add item to cart"})
			return
		}

		// Отправляем разные ответы в зависимости от действия
        if action == "added" {
            c.JSON(http.StatusOK, gin.H{"message": "Item added to favorite"})
        } else if action == "removed" {
            c.JSON(http.StatusOK, gin.H{"message": "Item removed from favorite"})
        }
	}
}

// Функция для добавления товара в корзину в репозитории
func (repo *LikedRepository) AddToLiked(userID, itemID uint) (string, error) {
    var likedItem models.LikedItem
    result := repo.db.Where("user_id = ? AND item_id = ?", userID, itemID).First(&likedItem)

    if result.Error == nil {
        // Если товар уже в избранных, удаляем его
        if err := repo.db.Unscoped().Delete(&likedItem).Error; err != nil {
            log.Println("Error deleting liked item:", err)
            return "", err
        }
        log.Println("Removed item from favorites:", likedItem)
        return "removed", nil
    } else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        // Если товара нет в избранных, создаем новую запись
        newLikedItem := models.LikedItem{
            UserID: userID,
            ItemID: itemID,
        }
        if err := repo.db.Create(&newLikedItem).Error; err != nil {
            log.Println("Error adding liked item:", err)
            return "", err
        }
        log.Println("Added item to favorites:", newLikedItem)
        return "added", nil
    } else {
        log.Println("Database error when checking for existing favorite item:", result.Error)
        return "", result.Error
    }
}


func (repo *LikedRepository) GetLikedItems(userID uint) ([]models.LikedItem, error) {
	var likedItems []models.LikedItem
	if err := repo.db.Preload("Item").Where("user_id = ?", userID).Find(&likedItems).Error; err != nil {
		log.Println("Error retrieving liked items:", err)
		return nil, err
	}
	log.Println("Retrieved liked items for user:", userID, likedItems)
	return likedItems, nil
}

func ViewLikedHandler(repo *LikedRepository, cartRepo *CartRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userID")
        if !exists {
            log.Println("Пользователь не аутентифицирован")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
            return
        }

        userIDUint, ok := userID.(uint)
        if !ok {
            log.Println("Неверный ID пользователя:", userID)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный ID пользователя"})
            return
        }

        // Получаем товары в избранных
        likedItems, err := repo.GetLikedItems(userIDUint)
        if err != nil {
            log.Println("Не удалось получить товары в избранных:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить товары в избранных"})
            return
        }

        // Создаем карту корзины
        cartItems, err := cartRepo.GetCartItems(userIDUint)
        if err != nil {
            log.Println("Не удалось получить товары в корзине:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить товары в корзине"})
            return
        }

        cartMap := make(map[uint]bool)
        for _, item := range cartItems {
            cartMap[item.ItemID] = true
        }

        // Получаем изображение, рейтинг и количество отзывов для каждого товара в избранных
        for i, likedItem := range likedItems {
            // Получаем изображение товара
            var productImage models.ItemImage
            if err := repo.db.Where("item_id = ?", likedItem.ItemID).First(&productImage).Error; err == nil {
                likedItems[i].Item.ImagePath = productImage.ImagePath
            }

            // Вычисление среднего рейтинга для товара
            var averageRating float64
            repo.db.Model(&models.Review{}).Where("item_id = ?", likedItem.ItemID).Select("AVG(rating)").Scan(&averageRating)
            likedItems[i].Item.AverageRating = averageRating

            // Считывание общего количества отзывов
            var reviewCount int64
            repo.db.Model(&models.Review{}).Where("item_id = ?", likedItem.ItemID).Count(&reviewCount)
            likedItems[i].Item.ReviewCount = reviewCount
        }

        // Подготовка данных для шаблона
        log.Println("Рендеринг страницы избранных для пользователя:", userIDUint)
        c.HTML(http.StatusOK, "liked.html", gin.H{
            "LikedItems": likedItems,
            "CartMap":    cartMap,  // Передаем карту корзины
        })
    }
}

