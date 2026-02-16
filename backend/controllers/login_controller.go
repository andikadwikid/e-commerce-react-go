package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend-commerce/database"
	"backend-commerce/helpers"
	"backend-commerce/models"
	"backend-commerce/structs"

)

func Login(c *gin.Context) {

	var req = structs.LoginRequest{}
	var user = models.User{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err, req),
		})
		return
	}

	if err := database.DB.Preload("Roles").Preload("Roles.Permissions").Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "User Not Found",
			Errors:  map[string]string{"email": "User not found or invalid credentials"},
		})
		return
	}

	if !helpers.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "Invalid Password",
			Errors:  map[string]string{"password": "Invalid password"},
		})
		return
	}

	token, err := helpers.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	permissionMap := helpers.GetPermissionMap(user.Roles)

	userResponse := structs.ToUserLoginResponse(user, permissionMap)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Login Success",
		Data: structs.LoginResponse{
			Token: token,
			User:  userResponse,
		},
	})
}
