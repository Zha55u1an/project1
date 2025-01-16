// itemhandlers.go

package handlers

import (
	"fmt"
	"go_project/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

type RecentlyViewedRepository struct {
	db *gorm.DB
}

func NewRecentlyViewedRepository(db *gorm.DB) *RecentlyViewedRepository {
	return &RecentlyViewedRepository{db: db}
}


func (repo *ItemRepository) GetAllItems(c *gin.Context) {
	var items []models.Item
	result := repo.db.Find(&items)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (repo *ItemRepository) GetItemByID(c *gin.Context, recentlyViewedRepo *RecentlyViewedRepository) {
    id := c.Param("id")
    var item models.Item
    var images []models.ItemImage
    var reviews []models.Review
    var averageRating float64
    var reviewCount int64

    // Находим товар по ID
    if err := repo.db.First(&item, id).Error; err != nil {
        c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Item not found"})
        return
    }

    // Добавляем товар в недавно просмотренные
    if userID, exists := c.Get("userID"); exists {
        err := recentlyViewedRepo.AddRecentlyViewed(userID.(uint), item.ID)
        if err != nil {
            log.Printf("Error adding recently viewed item: %v", err)
        }
    }

    // Находим все изображения, связанные с товаром
    repo.db.Where("item_id = ?", id).Find(&images)

    // Получаем отзывы для товара
    repo.db.Where("item_id = ?", id).Find(&reviews)

    // Вычисляем средний рейтинг
    repo.db.Model(&models.Review{}).Where("item_id = ?", id).Select("AVG(rating)").Scan(&averageRating)

    // Подсчитываем количество отзывов
    repo.db.Model(&models.Review{}).Where("item_id = ?", id).Count(&reviewCount)

    // Логика для вычисления количества звезд
    for i, review := range reviews {
        reviews[i].FilledStars = review.Rating         // Заполненные звезды (по рейтингу)
        reviews[i].RemainingStars = 5 - review.Rating // Оставшиеся звезды
    }

    // Передаем данные в шаблон
    c.HTML(http.StatusOK, "item.html", gin.H{
        "Item":          item,
        "Images":        images,
        "Reviews":       reviews,
        "AverageRating": averageRating,
        "ReviewCount":   reviewCount, // Количество отзывов
    })
}




func (repo *RecentlyViewedRepository) AddRecentlyViewed(userID uint, itemID uint) error {
    // Убедимся, что запись не дублируется
    var existingItem models.RecentlyViewedItem
    if err := repo.db.Where("user_id = ? AND item_id = ?", userID, itemID).First(&existingItem).Error; err == nil {
        // Если запись уже существует, обновляем время просмотра
        existingItem.ViewTime = time.Now()
        if err := repo.db.Save(&existingItem).Error; err != nil {
            log.Printf("Error updating recently viewed item: %v", err)
            return err
        }
        return nil
    }

    // Если запись не существует, создаем новую
    newRecentlyViewed := models.RecentlyViewedItem{
        UserID:   userID,
        ItemID:   itemID,
        ViewTime: time.Now(),
    }
    if err := repo.db.Create(&newRecentlyViewed).Error; err != nil {
        log.Printf("Error adding recently viewed item: %v", err)
        return err
    }

    return nil
}





func (repo *ItemRepository) CreateItem(c *gin.Context) {
	var item models.Item
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	result := repo.db.Where("name = ?", item.Category.Name).First(&category)

	if category.ID == 0 {
		repo.db.Create(category)
		fmt.Println("my category =", category)
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	item.Category = &category

	result = repo.db.Create(&item)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (repo *ItemRepository) UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := repo.db.Model(&models.Item{}).Where("id = ?", id).Updates(&item)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	repo.db.Where("id=?", id).First(&item)
	c.JSON(http.StatusOK, item)
}

func (repo *ItemRepository) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	result := repo.db.Delete(&models.Item{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
