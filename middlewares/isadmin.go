// PATH: go-auth/middlewares/isAuthorized.go

package middlewares

import (

	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
		role, exists := c.Get("role") 
        
        if !exists {
            c.JSON(401, gin.H{"error": "user is not an admin"})
            c.Abort()
            return
        }

        if role != "admin" {
            c.JSON(401, gin.H{"error": "user is not an admin"})
            c.Abort()
            return
        }
        c.Next()
    }
}
