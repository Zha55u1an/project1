package handlers

import (
	"fmt"
	"go_project/internal/models"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (repo *ItemRepository) SellerDashboard(c *gin.Context) {
	// Получаем ID текущего продавца из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Статистика товаров
	var totalItems int64
	repo.db.Model(&models.Item{}).Where("seller_id = ?", userID).Count(&totalItems)

	// Статистика заказов
	var totalOrders int64
	repo.db.Model(&models.Order{}).Where("seller_id = ?", userID).Count(&totalOrders)

	// Передаем статистику в шаблон
	c.HTML(http.StatusOK, "seller-dashboard.html", gin.H{
		"TotalItems":   totalItems,
		"TotalOrders":  totalOrders,
		"UserID":       userID,
	})
}


func (repo *ItemRepository) SellerGetAllItems(c *gin.Context) {
	// Получаем ID текущего пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var items []models.Item
	// Загружаем все товары для текущего продавца
	result := repo.db.Where("seller_id = ?", userID).Find(&items)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Присваиваем путь первого изображения каждому товару
	for i, item := range items {
		// Для каждого товара получаем его изображения (первое)
		var itemImage models.ItemImage
		if err := repo.db.Where("item_id = ?", item.ID).First(&itemImage).Error; err == nil {
			// Присваиваем путь первого изображения товару
			items[i].ImagePath = itemImage.ImagePath
		}
	}

	// Отправляем товары с изображениями на страницу
	c.HTML(http.StatusOK, "seller-products.html", gin.H{"items": items})
}

func (repo *ItemRepository) SellerCreateItem(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		log.Println("Получен GET-запрос для отображения формы создания товара")
		c.HTML(http.StatusOK, "create_item.html", gin.H{})
		return
	}

	var item models.Item

	// Извлечение ID продавца
	sellerID, exists := c.Get("userID")
	if !exists {
		log.Println("Ошибка: ID продавца не найден в контексте")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не удалось определить продавца"})
		return
	}
	item.SellerID = sellerID.(uint)
	log.Printf("ID продавца: %d\n", item.SellerID)

	// Извлечение текстовых данных из формы
	item.Name = c.PostForm("name")
	item.Description = c.PostForm("description")
	log.Printf("Извлеченные данные товара: Название - %s, Описание - %s", item.Name, item.Description)

	// Извлечение и конвертация цены
	priceStr := c.PostForm("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Println("Ошибка при конвертации цены:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат цены"})
		return
	}
	item.Price = price
	log.Printf("Цена товара: %.2f", item.Price)

	// Проверка доступности
	isAvailableStr := c.PostForm("is_available")
	item.IsAvailable = isAvailableStr == "Yes"
	log.Printf("Товар доступен для продажи: %v", item.IsAvailable)

	// Получаем значение категории отдельно
	itemCategoryName := c.PostForm("category")
	if itemCategoryName != "" {
		item.Category = &models.Category{Name: itemCategoryName}
		log.Printf("Категория товара: %s", item.Category.Name)
	} else {
		log.Println("Ошибка: категория не указана")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Категория не указана"})
		return
	}

	// Проверка категории и создание новой при необходимости
	var category models.Category
	result := repo.db.Where("name = ?", item.Category.Name).First(&category)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Println("Ошибка при проверке категории в базе данных:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if category.ID == 0 {
		log.Printf("Категория не найдена, создаем новую категорию: %s", item.Category.Name)
		newCategory := models.Category{Name: item.Category.Name}
		if err := repo.db.Create(&newCategory).Error; err != nil {
			log.Println("Ошибка при создании новой категории:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		category = newCategory
		log.Printf("Новая категория успешно создана: %s", category.Name)
	}
	item.Category = &category

	// Создаем товар в базе данных
	if err := repo.db.Create(&item).Error; err != nil {
		log.Println("Ошибка при создании товара в базе данных:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Товар успешно создан с ID: %d", item.ID)

	// Загрузка до 10 изображений
	files := c.Request.MultipartForm.File["images"]
	log.Printf("Загружаем %d изображений", len(files))
	if len(files) > 10 {
		log.Println("Ошибка: превышен лимит изображений")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вы можете загрузить максимум 10 изображений"})
		return
	}

	for _, file := range files {
		fileHeader, err := file.Open()
		if err != nil {
			log.Println("Ошибка при открытии файла:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла"})
			return
		}
		defer fileHeader.Close()

		// Генерация уникального имени файла
		uniqueId := uuid.New()
		filename := strings.Replace(uniqueId.String(), "-", "", -1)
		fileExt := filepath.Ext(file.Filename)
		imageName := fmt.Sprintf("%s%s", filename, fileExt)
		filePath := fmt.Sprintf("./uploads/%s", imageName)

		// Сохраняем файл на сервере
		outFile, err := os.Create(filePath)
		if err != nil {
			log.Println("Ошибка при сохранении изображения:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения изображения"})
			return
		}
		_, err = io.Copy(outFile, fileHeader)
		if err != nil {
			log.Println("Ошибка при копировании файла:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения изображения"})
			return
		}

		// Сохраняем данные изображения в базе данных
		image := models.ItemImage{
			ItemID:    item.ID,
			ImagePath: fmt.Sprintf("http://localhost:8080/uploads/%s", imageName),
		}
		log.Printf("Сохраняем изображение в БД: %+v", image)
		if err := repo.db.Create(&image).Error; err != nil {
			log.Println("Ошибка при сохранении данных изображения в БД:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения изображения в БД"})
			return
		}
		log.Printf("Изображение успешно сохранено: %s", image.ImagePath)
	}

	log.Printf("Товар успешно создан с изображениями: %v", item)
	c.JSON(http.StatusCreated, gin.H{"item": item, "message": "Item created successfully!"})
}

func (repo *ItemRepository) SellerEditItem(c *gin.Context) {
	// Получаем ID товара из параметров маршрута
	itemID := c.Param("id")

	var item models.Item
	// Ищем товар по его ID, а также подгружаем категорию товара
	if err := repo.db.Preload("Category").First(&item, itemID).Error; err != nil {
		log.Println("Ошибка при поиске товара:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
		return
	}

	// Получаем все изображения для этого товара
	var itemImages []models.ItemImage
	repo.db.Where("item_id = ?", item.ID).Find(&itemImages)

	// Проверяем, является ли текущий пользователь владельцем товара
	userID, exists := c.Get("userID")
	if !exists || item.SellerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете редактировать этот товар"})
		return
	}

	// Отправляем данные товара и его изображения на страницу редактирования
	c.HTML(http.StatusOK, "edit-item.html", gin.H{
		"item":       item,
		"images":     itemImages, // Передаем список изображений
	})
}




func (repo *ItemRepository) SellerUpdateItem(c *gin.Context) {
	// Получаем ID товара из параметров маршрута
	itemID := c.Param("id")

	var item models.Item
	// Ищем товар по его ID
	if err := repo.db.Preload("Category").First(&item, itemID).Error; err != nil {
		log.Println("Ошибка при поиске товара:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
		return
	}

	// Проверяем, является ли текущий пользователь владельцем товара
	userID, exists := c.Get("userID")
	if !exists || item.SellerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете редактировать этот товар"})
		return
	}

	// Извлекаем обновленные данные из формы
	item.Name = c.PostForm("name")
	item.Description = c.PostForm("description")
	priceStr := c.PostForm("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Println("Ошибка при конвертации цены:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат цены"})
		return
	}
	item.Price = price

	isAvailableStr := c.PostForm("is_available")
	item.IsAvailable = isAvailableStr == "Yes"

	itemCategoryName := c.PostForm("category")

	// Обновляем категорию товара
	var category models.Category
	if itemCategoryName != "" {
		result := repo.db.Where("name = ?", itemCategoryName).First(&category)
		if result.Error != nil {
			// Если категория не найдена, создаем новую
			log.Printf("Категория %s не найдена, создаем новую", itemCategoryName)
			newCategory := models.Category{Name: itemCategoryName}
			if err := repo.db.Create(&newCategory).Error; err != nil {
				log.Println("Ошибка при создании новой категории:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			category = newCategory
		}
		// Присваиваем товару категорию
		item.Category = &category
	}

	// Обновляем товар в базе данных
	if err := repo.db.Save(&item).Error; err != nil {
		log.Println("Ошибка при обновлении товара в базе данных:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновление изображений (если есть новые изображения)
	files := c.Request.MultipartForm.File["images"]
	if len(files) > 0 {
		// Удаляем старые изображения
		var itemImages []models.ItemImage
		repo.db.Where("item_id = ?", item.ID).Delete(&itemImages)

		// Загружаем новые изображения
		for _, file := range files {
			fileHeader, err := file.Open()
			if err != nil {
				log.Println("Ошибка при открытии файла:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла"})
				return
			}
			defer fileHeader.Close()

			// Генерация уникального имени файла
			uniqueId := uuid.New()
			filename := strings.Replace(uniqueId.String(), "-", "", -1)
			fileExt := filepath.Ext(file.Filename)
			imageName := fmt.Sprintf("%s%s", filename, fileExt)
			filePath := fmt.Sprintf("./uploads/%s", imageName)

			// Сохраняем файл на сервере
			outFile, err := os.Create(filePath)
			if err != nil {
				log.Println("Ошибка при сохранении изображения:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения изображения"})
				return
			}
			_, err = io.Copy(outFile, fileHeader)
			if err != nil {
				log.Println("Ошибка при копировании файла:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения изображения"})
				return
			}

			// Сохраняем данные изображения в базе данных
			image := models.ItemImage{
				ItemID:    item.ID,
				ImagePath: fmt.Sprintf("http://localhost:8080/uploads/%s", imageName),
			}
			if err := repo.db.Create(&image).Error; err != nil {
				log.Println("Ошибка при сохранении изображения в БД:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения изображения в БД"})
				return
			}
		}
	}

	// После успешного обновления товара и его изображений, перенаправляем на страницу с продуктами
	c.HTML(http.StatusOK, "seller-products.html", gin.H{})
}


func (repo *ItemRepository) SellerDeleteItem(c *gin.Context) {
	// Получаем ID товара из параметров маршрута
	itemID := c.Param("id")

	// Находим товар
	var item models.Item
	if err := repo.db.First(&item, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Проверяем, является ли текущий пользователь владельцем товара
	userID, exists := c.Get("userID")
	if !exists || item.SellerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this item"})
		return
	}

	// Удаляем все изображения товара
	var itemImages []models.ItemImage
	repo.db.Where("item_id = ?", item.ID).Find(&itemImages)
	for _, image := range itemImages {
		if err := repo.db.Delete(&image).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting images"})
			return
		}
	}

	// Удаляем товар
	if err := repo.db.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
