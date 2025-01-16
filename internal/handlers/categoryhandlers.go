// categoryhandlers.go

package handlers

import (
    "github.com/gin-gonic/gin"
    "go_project/internal/models"
    "gorm.io/gorm"
    "net/http"
)

type CategoryRepository struct {
    db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
    return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetAllCategories(c *gin.Context) {
    var categories []models.Category
    result := repo.db.Preload("Products").Find(&categories)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, categories)
}

func (repo *CategoryRepository) GetCategoryByID(c *gin.Context) {
    id := c.Param("id")
    var category []models.Category
    result := repo.db.First(&category, id)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
        return
    }
    c.JSON(http.StatusOK, category)
}

func (repo *CategoryRepository) CreateCategory(c *gin.Context) {
    var category models.Category
    if err := c.BindJSON(&category); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result := repo.db.Create(&category)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusCreated, category)
}

func (repo *CategoryRepository) UpdateCategory(c *gin.Context) {
    id := c.Param("id")
    var category models.Category
    if err := c.BindJSON(&category); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result := repo.db.Model(&models.Category{}).Where("id = ?", id).Updates(&category)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    repo.db.Where("id=?", id).First(&category)
    c.JSON(http.StatusOK, category)
}

func (repo *CategoryRepository) DeleteCategory(c *gin.Context) {
    id := c.Param("id")
    result := repo.db.Delete(&models.Category{}, id)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
