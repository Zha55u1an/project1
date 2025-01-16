package handlers

import (
	"go_project/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomePageData struct {
    User     string
    Products []models.Item // Здесь предполагается, что модель Item уже существует
}

func (repo *ItemRepository) HomePage(cartRepo *CartRepository, likedRepo *LikedRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        var products []models.Item

        // Извлекаем все доступные товары
        result := repo.db.Where("is_available = ?", true).Find(&products)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve products"})
            return
        }

        // Добавляем изображения для товаров и вычисляем рейтинг
        for i, product := range products {
            var productImage models.ItemImage
            if err := repo.db.Where("item_id = ?", product.ID).First(&productImage).Error; err == nil {
                products[i].ImagePath = productImage.ImagePath
            }

            // Вычисление среднего рейтинга для товара
            var averageRating float64
            repo.db.Model(&models.Review{}).Where("item_id = ?", product.ID).Select("AVG(rating)").Scan(&averageRating)
            products[i].AverageRating = averageRating

            // Получаем количество отзывов для товара
            var reviewCount int64
            repo.db.Model(&models.Review{}).Where("item_id = ?", product.ID).Count(&reviewCount)
            products[i].ReviewCount = reviewCount
        }

        // Получаем ID текущего пользователя
        var userID uint
        if userIDValue, exists := c.Get("userID"); exists {
            userID = userIDValue.(uint)
        }

        // Проверяем, какие товары в избранных у текущего пользователя
        likedItems := []models.LikedItem{}
        likedRepo.db.Where("user_id = ?", userID).Find(&likedItems)

        // Создаем карту, где ключ — `item_id`, а значение — true, если товар в избранных
        likedMap := make(map[uint]bool)
        for _, likedItem := range likedItems {
            likedMap[likedItem.ItemID] = true
        }

        // Получаем список товаров в корзине
        var cartItems []models.CartItem
        cartRepo.db.Where("user_id = ?", userID).Find(&cartItems)

        // Создаем карту товаров в корзине для быстрого доступа
        cartMap := make(map[uint]bool)
        for _, item := range cartItems {
            cartMap[item.ItemID] = true
        }

        // Подготовка данных для шаблона
        data := gin.H{
            "User":     userID,
            "Products": products,
            "LikedMap": likedMap, // Передаем карту избранных
            "CartMap":  cartMap,
        }

        // Отправляем данные на home.html
        c.HTML(http.StatusOK, "home.html", data)
    }
}







