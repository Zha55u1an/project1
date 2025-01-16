// PATH: go-auth/middlewares/isAuthorized.go

package middlewares

import (
	"fmt"
	"go_project/internal/models"
	"go_project/pkg/db"
	"go_project/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из cookie
		cookie, err := c.Cookie("token")
		fmt.Println("Cookie:", cookie)
		if err != nil {
			fmt.Println("Error retrieving cookie:", err)
			// Перенаправляем на страницу входа
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Парсим токен и проверяем его
		claims, err := utils.ParseToken(cookie)
		if err != nil {
			fmt.Println("Error parsing token:", err)
			// Перенаправляем на страницу входа
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Проверяем, существует ли пользователь в базе данных
		var user models.User
		if err := db.DB.Where("email = ?", claims.Subject).First(&user).Error; err != nil {
			fmt.Println("Error finding user:", err)
			// Перенаправляем на страницу входа
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Устанавливаем данные пользователя в контекст
		c.Set("role", claims.Role)
		c.Set("userID", user.ID)

		// Передаём управление следующему обработчику
		c.Next()
	}
}

