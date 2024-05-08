// authentication.go

package handlers

import (
	"go_project/internal/models"
	"go_project/pkg/db"
	"go_project/pkg/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Login обрабатывает POST запрос для аутентификации пользователя
// PATH: go-auth/controllers/auth.go
var jwtKey = []byte("jwtkey_go_project")

func Login(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	db.DB.Where("username = ?", user.Username).First(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	errHash := utils.CheckPasswordHash(user.Password, existingUser.Password)

	if !errHash {
		c.JSON(400, gin.H{"error": "invalid password"})
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)

	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Username,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged in"})
}

// PATH: go-auth/controllers/auth.go

func Signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	db.DB.Where("username = ?", user.Username).First(&existingUser)

	if existingUser.ID != 0 {
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}

	var errHash error
	user.Password, errHash = utils.PasswordHash(user.Password)

	if errHash != nil {
		c.JSON(500, gin.H{"error": "could not generate password hash"})
		return
	}

	db.DB.Create(&user)

	c.JSON(200, gin.H{"success": "user created"})
}

// Logout обрабатывает POST запрос для выхода пользователя из системы
// PATH: go-auth/controllers/auth.go

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged out"})
}
