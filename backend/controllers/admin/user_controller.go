package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend-commerce/database"
	"backend-commerce/helpers"
	"backend-commerce/models"
	"backend-commerce/structs"
)

func FindUsers(c *gin.Context) {
	var users []models.User
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Model(&models.User{})
	if search != "" {
		query = query.Where("name LIKE ? OR username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to count users",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	if err := query.Preload("Roles").Order("id desc").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch users",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	var data []structs.UserDetailResponse
	for _, u := range users {
		var userRoles []structs.RoleResponse
		for _, r := range u.Roles {
			userRoles = append(userRoles, structs.RoleResponse{
				Id:   r.Id,
				Name: r.Name,
			})
		}

		data = append(data, structs.UserDetailResponse{
			Id:       u.Id,
			Name:     u.Name,
			Username: u.Username,
			Email:    u.Email,
			Roles:    userRoles,
		})
	}

	helpers.PaginateResponse(c, data, total, page, limit, baseURL, search, "List Data Users")
}

func CreateUser(c *gin.Context) {
	var request structs.UserCreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Failed",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	hashedPassword, err := helpers.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to hash password",
		})
		return
	}

	user := models.User{
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
	}

	var roles []models.Role
	if len(request.RoleIDs) > 0 {
		database.DB.Where("id IN ?", request.RoleIDs).Find(&roles)
	}
	user.Roles = roles

	if err := database.DB.Create(&user).Error; err != nil {
		if helpers.IsDuplicateEntryError(err) {
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Success: false,
				Message: "Create User Failed",
				Errors:  helpers.TranslateErrorMessage(err, nil),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "User Created Successfully",
		Data: gin.H{
			"id":       user.Id,
			"name":     user.Name,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	var request structs.UserUpdateRequest

	if err := database.DB.Preload("Roles").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User Not Found",
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Failed",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	user.Name = request.Name
	user.Username = request.Username
	user.Email = request.Email

	if request.Password != "" {
		hashed, err := helpers.HashPassword(request.Password)
		if err == nil {
			user.Password = hashed
		}
	}

	var newRoles []models.Role
	if len(request.RoleIDs) > 0 {
		database.DB.Where("id IN ?", request.RoleIDs).Find(&newRoles)
	}
	if err := database.DB.Model(&user).Association("Roles").Replace(newRoles); err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update user roles",
		})
		return
	}

	// 5. Simpan User
	if err := database.DB.Save(&user).Error; err != nil {
		if helpers.IsDuplicateEntryError(err) {
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Success: false,
				Message: "Update User Failed (Duplicate Data)",
				Errors:  helpers.TranslateErrorMessage(err, nil),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update user",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User Updated Successfully",
	})
}

func GetUserDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.User

	if err := database.DB.Preload("Roles").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User Not Found",
		})
		return
	}

	var userRoles []structs.RoleResponse
	for _, r := range user.Roles {
		userRoles = append(userRoles, structs.RoleResponse{
			Id:   r.Id,
			Name: r.Name,
		})
	}
	res := structs.UserDetailResponse{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Roles:    userRoles,
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User Detail",
		Data:    res,
	})
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User Not Found",
		})
		return
	}

	if err := database.DB.Select("Roles").Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete user",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User Deleted Successfully",
	})
}
