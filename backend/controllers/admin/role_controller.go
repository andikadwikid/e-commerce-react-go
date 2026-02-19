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

func GetAllRoles(c *gin.Context) {
	var roles []models.Role
	if err := database.DB.Preload("Permissions").Order("name asc").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch roles",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "All Roles List",
		Data:    roles,
	})
}

// FindRoles mengambil list role dengan pagination & search
func FindRoles(c *gin.Context) {
	var roles []models.Role
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Model(&models.Role{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to count roles",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	// Preload permissions agar terlihat di list (optional, bisa dihilangkan jika berat)
	if err := query.Preload("Permissions").Order("id desc").Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch roles",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	helpers.PaginateResponse(c, roles, total, page, limit, baseURL, search, "List Data Roles")
}

func CreateRole(c *gin.Context) {
	var request structs.RoleCreateRequest

	// 1. Validasi Input
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Failed",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	// 2. Siapkan Model
	role := models.Role{
		Name: request.Name,
	}

	// 3. Cari permissions berdasarkan ID yang dikirim
	var permissions []models.Permission
	if len(request.PermissionIDs) > 0 {
		database.DB.Where("id IN ?", request.PermissionIDs).Find(&permissions)
	}
	role.Permissions = permissions

	// 4. Simpan Role + Relasi Permissions
	if err := database.DB.Create(&role).Error; err != nil {
		if helpers.IsDuplicateEntryError(err) {
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Success: false,
				Message: "Create Role Failed",
				Errors:  helpers.TranslateErrorMessage(err, nil),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create role",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Role Created Successfully",
		Data:    role,
	})
}

func GetRoleDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role models.Role

	if err := database.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Role Detail",
		Data:    role,
	})
}

func UpdateRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role models.Role
	var request structs.RoleUpdateRequest

	// 1. Cek Role
	if err := database.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role Not Found",
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

	role.Name = request.Name

	var newPermissions []models.Permission
	if len(request.PermissionIDs) > 0 {
		database.DB.Where("id IN ?", request.PermissionIDs).Find(&newPermissions)
	}

	if err := database.DB.Model(&role).Association("Permissions").Replace(newPermissions); err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update role permissions",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	if err := database.DB.Save(&role).Error; err != nil {
		if helpers.IsDuplicateEntryError(err) {
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Success: false,
				Message: "Update Role Failed",
				Errors:  helpers.TranslateErrorMessage(err, nil),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update role",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}
	database.DB.Preload("Permissions").First(&role, id)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Role Updated Successfully",
		Data:    role,
	})
}

func DeleteRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role models.Role

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role Not Found",
		})
		return
	}

	if err := database.DB.Select("Permissions").Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete role",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Role Deleted Successfully",
	})
}
