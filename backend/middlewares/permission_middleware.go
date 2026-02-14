package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend-commerce/database"
	"backend-commerce/models"
)

// Middleware untuk cek apakah user memiliki permission tertentu
func Permission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil username dari context (disimpan oleh middleware Auth)
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var user models.User
		err := database.DB.
			Preload("Roles.Permissions").
			Where("username = ?", username).
			First(&user).Error

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Cek apakah user memiliki permission yang diminta
		// Loop level 1: Role
		for _, role := range user.Roles {
			// Loop level 2: Permission di dalam Role tersebut
			for _, perm := range role.Permissions {
				if perm.Name == permissionName {
					c.Next()
					return
				}
			}
		}

		// Jika sampai sini berarti tidak ada permission yang cocok
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden - You don't have permission to access this resource"})
		c.Abort()
	}
}
