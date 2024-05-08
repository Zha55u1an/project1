// PATH: go-auth/middlewares/isAuthorized.go

package middlewares

import (
	"fmt"
	"go_project/internal/models"
	"go_project/pkg/db"
	"go_project/pkg/utils"

	"github.com/gin-gonic/gin"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		fmt.Println("Cookie", cookie)
		if err != nil {
			fmt.Println(err)
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(cookie)

		if err != nil {
			fmt.Println(err)
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		var user models.User
		db.DB.Where("username = ?", claims.Subject).First(&user)

		c.Set("role", claims.Role)
		c.Set("userID", user.ID)
		c.Next()
	}
}
