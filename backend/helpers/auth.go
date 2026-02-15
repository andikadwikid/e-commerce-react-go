package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"backend-commerce/database"
	"backend-commerce/models"

)

func GetAuthUserID(c *gin.Context) (uint, error) {
	username, exists := c.Get("username")
	if !exists {
		return 0, errors.New("user not found in context")
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}
	return user.Id, nil
}
