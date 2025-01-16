package middlewares

import (
	"github.com/gin-gonic/gin"
)

func SellerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role") 
        
        if !exists {
            c.JSON(401, gin.H{"error": "user is not seller"})
            c.Abort()
            return
        }

        if role != "seller" {
            c.JSON(401, gin.H{"error": "user is not seller"})
            c.Abort()
            return
        }
        c.Next()
    }
}
