// itemhandlers.go

package handlers

import (
	"fmt"
	"go_project/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// GetAllItems handles GET request to fetch all items
func (repo *ItemRepository) GetAllItems(c *gin.Context) {
	var items []models.Item
	result := repo.db.Find(&items)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetItemByID handles GET request to fetch an item by ID
func (repo *ItemRepository) GetItemByID(c *gin.Context) {
	id := c.Param("id")
	var item models.Item
	result := repo.db.First(&item, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateItem handles POST request to create a new item
func (repo *ItemRepository) CreateItem(c *gin.Context) {
	var item models.Item
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск категории по имени (или создание, если она не существует)
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

	// Связывание товара с найденной/созданной категорией
	item.Category = &category

	// Создание товара (включая связанную категорию)
	result = repo.db.Create(&item)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdateItem handles PUT request to update an existing item
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

// DeleteItem handles DELETE request to delete an item by ID
func (repo *ItemRepository) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	result := repo.db.Delete(&models.Item{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
