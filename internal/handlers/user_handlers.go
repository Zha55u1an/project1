// user_handler.go

package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_project/internal/models"
	"gorm.io/gorm"
	"net/http"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAllUsers обрабатывает GET запрос для получения всех пользователей
func (repo *UserRepository) GetAllUsers(c *gin.Context) {
	fmt.Println("ASDSD")
	var users []models.User
	result := repo.db.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSONP(http.StatusOK, users)
}

// GetUserByID обрабатывает GET запрос для получения пользователя по ID
func (repo *UserRepository) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	result := repo.db.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser обрабатывает POST запрос для создания нового пользователя

// UpdateUser обрабатывает PUT запрос для обновления информации о пользователе
func (repo *UserRepository) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := repo.db.Model(&models.User{}).Where("id = ?", id).Updates(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser обрабатывает DELETE запрос для удаления пользователя по ID
func (repo *UserRepository) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	result := repo.db.Delete(&models.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
