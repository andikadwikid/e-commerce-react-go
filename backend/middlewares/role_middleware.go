package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend-commerce/database"
	"backend-commerce/models"
)

// Middleware untuk cek apakah user memiliki role tertentu
func Role(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var user models.User

		// Load user lengkap dengan relasi roles
		err := database.DB.Preload("Roles").Where("username = ?", username).First(&user).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
			c.Abort()
			return
		}

		// Cek apakah user memiliki role yang diminta
		hasRole := false
		for _, role := range user.Roles {
			if role.Name == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": http.StatusText(http.StatusForbidden)})
			c.Abort()
			return
		}

		c.Next()
	}
}
