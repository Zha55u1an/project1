package handlers

import (
	"fmt"
	"go_project/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ItemController struct {
    repo *ItemRepository
}

func NewItemController(repo *ItemRepository) *ItemController {
    return &ItemController{repo: repo}
}

func (repo *ItemRepository) SearchItems(query string, categoryID uint, maxPrice float64) ([]models.Item, error) {
    var items []models.Item

    // Строим запрос с учетом переданных параметров
    dbQuery := repo.db.Preload("Category") // Используем Preload для загрузки данных о категории товара

    // Если query задан, ищем товары по имени
    if query != "" {
        dbQuery = dbQuery.Where("name LIKE ?", "%"+query+"%")
    }

    // Если categoryID задан, добавляем фильтрацию по категории
    if categoryID > 0 {
        dbQuery = dbQuery.Where("category_id = ?", categoryID)
    }

    // Если maxPrice задан, фильтруем товары по цене
    if maxPrice > 0 {
        dbQuery = dbQuery.Where("price <= ?", maxPrice)
    }

    // Выполняем запрос
    if err := dbQuery.Find(&items).Error; err != nil {
        return nil, fmt.Errorf("Error occurred while querying items: %v", err)
    }

    return items, nil
}

// Обработчик поиска товаров
func (ctrl *ItemController) HandleSearchItems(c *gin.Context) {
    query := c.DefaultQuery("query", "") // Поиск по имени товара
    categoryID := c.DefaultQuery("category_id", "0") // Поиск по категории (0 по умолчанию для всех категорий)
    maxPrice, _ := strconv.ParseFloat(c.DefaultQuery("max_price", "0"), 64) // Поиск по максимальной цене

    // Преобразуем categoryID в uint
    catID, err := strconv.Atoi(categoryID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
        return
    }

    // Получаем товары из репозитория
    items, err := ctrl.repo.SearchItems(query, uint(catID), maxPrice)
    totalCount := len(items)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Получаем информацию о корзине и избранных товарах для текущего пользователя
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    userIDUint, ok := userID.(uint)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
        return
    }

    // Получаем товары в корзине
    var cartItems []models.CartItem
    ctrl.repo.db.Where("user_id = ?", userIDUint).Find(&cartItems)
    cartMap := make(map[uint]bool)
    for _, item := range cartItems {
        cartMap[item.ItemID] = true
    }

    // Получаем избранные товары
    var likedItems []models.LikedItem
    ctrl.repo.db.Where("user_id = ?", userIDUint).Find(&likedItems)
    likedMap := make(map[uint]bool)
    for _, likedItem := range likedItems {
        likedMap[likedItem.ItemID] = true
    }

    // Добавляем рейтинг, количество отзывов и изображения для каждого товара
    for i := range items {
        // Получаем изображение товара
        var productImage models.ItemImage
        if err := ctrl.repo.db.Where("item_id = ?", items[i].ID).First(&productImage).Error; err == nil {
            items[i].ImagePath = productImage.ImagePath
        }

        // Вычисление среднего рейтинга для товара
        var averageRating float64
        ctrl.repo.db.Model(&models.Review{}).Where("item_id = ?", items[i].ID).Select("AVG(rating)").Scan(&averageRating)
        items[i].AverageRating = averageRating

        // Получаем количество отзывов для товара
        var reviewCount int64
        ctrl.repo.db.Model(&models.Review{}).Where("item_id = ?", items[i].ID).Count(&reviewCount)
        items[i].ReviewCount = reviewCount
    }

    // Подготовка данных для шаблона
    data := gin.H{
        "Query":       query,
        "Items":       items,
        "TotalCount":  totalCount,
        "LikedMap":    likedMap,  // Передаем карту избранных товаров
        "CartMap":     cartMap,   // Передаем карту товаров в корзине
    }

    // Отправляем данные на страницу
    c.HTML(http.StatusOK, "search_results.html", data)
}

