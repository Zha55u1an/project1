// authentication.go

// package handlers
// import (
// 	"github.com/gin-gonic/gin"
// 	"go_project/internal/models"
// 	"go_project/pkg"
// 	"golang.org/x/crypto/bcrypt"
// 	"go_project/pkg/utils"
// 	"gorm.io/gorm"
// 	"net/http"
// )

// Login обрабатывает POST запрос для аутентификации пользователя
// func Login(c *gin.Context) {
// 	var loginData models.LoginData
// 	if err := c.BindJSON(&loginData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var user models.User
// 	result := pkg.DB.Where("username = ?", loginData.Username).First(&user)
// 	if result.Error != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
// 		return
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password)); err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
// 		return
// 	}
//
//	// TODO: Generate JWT token and return it to the client
// }
//
//// Logout обрабатывает POST запрос для выхода пользователя из системы
// func Logout(c *gin.Context) {
// 	// TODO: Invalidate JWT token or clear session data
// 	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
// }