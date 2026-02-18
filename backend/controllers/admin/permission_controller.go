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

// Digunakan untuk dropdown list atau checklist saat membuat Role
func GetAllPermissions(c *gin.Context) {
	var permissions []models.Permission
	if err := database.DB.Order("name asc").Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch permissions",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "All Permissions List",
		Data:    permissions,
	})
}

// FindPermissions mengambil list permission dengan pagination & search.
// Flow besarnya:
// 1) Ambil parameter pencarian dan pagination dari request
// 2) Bentuk query dasar ke tabel permissions
// 3) Jika ada kata kunci pencarian, filter data berdasarkan nama
// 4) Hitung total data yang cocok dengan filter
// 5) Ambil data permissions sesuai page, limit, dan urutan
// 6) Bungkus data + info pagination ke dalam response standar
func FindPermissions(c *gin.Context) {
	// LANGKAH 1: Ambil parameter search dan pagination (page, limit, offset) dari request
	search, page, limit, offset := helpers.GetPaginationParams(c)

	// Base URL dipakai untuk membentuk link pagination di response
	baseURL := helpers.BuildBaseURL(c)

	// LANGKAH 2: Bentuk query dasar ke tabel permissions
	var permissions []models.Permission
	var total int64
	query := database.DB.Model(&models.Permission{})

	// LANGKAH 3: Jika ada kata kunci search, tambahkan filter ke query
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// LANGKAH 4: Hitung total data yang cocok dengan filter (tanpa limit & offset)
	if err := query.Count(&total).Error; err != nil {
		// Jika gagal menghitung total, balas dengan error 500
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to count permissions",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	// LANGKAH 5: Ambil data permissions untuk halaman tertentu (dengan limit, offset, dan urutan)
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&permissions).Error; err != nil {
		// Jika gagal mengambil data, balas dengan error 500
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch permissions",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	// LANGKAH 6: Kirim response dengan format pagination standar
	helpers.PaginateResponse(c, permissions, total, page, limit, baseURL, search, "List Data Permissions")
}

// CreatePermission menambahkan permission baru.
// Flow besarnya:
// 1) Ambil dan validasi data permission dari body request
// 2) Bentuk struct permission dari data yang sudah valid
// 3) Simpan permission baru ke database
// 4) Jika terjadi error duplikasi, balas dengan status conflict (409)
// 5) Jika error lain, balas dengan status 500
// 6) Jika sukses, balas dengan status 201 dan data permission baru
func CreatePermission(c *gin.Context) {
	// LANGKAH 1: Ambil dan validasi data permission dari body request
	var request structs.PermissionCreateRequest

	// Bind JSON ke struct request sekaligus validasi sesuai aturan di struct
	if err := c.ShouldBindJSON(&request); err != nil {
		// Jika data tidak valid, balas dengan status 422 dan detail error
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Failed",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	// LANGKAH 2: Bentuk struct permission dari data yang sudah valid
	permission := models.Permission{
		Name: request.Name,
	}

	// LANGKAH 3: Simpan permission baru ke database
	if err := database.DB.Create(&permission).Error; err != nil {
		// LANGKAH 4: Jika error karena duplikasi (nama permission sudah ada)
		if helpers.IsDuplicateEntryError(err) {
			// Balas dengan status 409 (Conflict)
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Success: false,
				Message: "Create Permission Failed",
				Errors:  helpers.TranslateErrorMessage(err, nil),
			})
			return
		}

		// LANGKAH 5: Jika error lain (bukan duplikasi), anggap error server
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create permission",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	// LANGKAH 6: Jika semua berhasil, kirim response sukses 201 dengan data permission
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Permission Created Successfully",
		Data:    permission,
	})
}

func GetPermissionDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var permission models.Permission

	if err := database.DB.First(&permission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Permission Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Permission Detail",
		Data:    permission,
	})
}

// UpdatePermission melakukan update permission
func UpdatePermission(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var permission models.Permission
	var request structs.PermissionCreateRequest

	// 1. Cek Data
	if err := database.DB.First(&permission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Permission Not Found",
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

	permission.Name = request.Name
	if err := database.DB.Save(&permission).Error; err != nil {
		if helpers.IsDuplicateEntryError(err) {
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Success: false,
				Message: "Update Permission Failed",
				Errors:  helpers.TranslateErrorMessage(err, nil),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update permission",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Permission Updated Successfully",
		Data:    permission,
	})
}

// DeletePermission menghapus permission berdasarkan ID
func DeletePermission(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var permission models.Permission

	if err := database.DB.First(&permission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Permission Not Found",
		})
		return
	}

	if err := database.DB.Delete(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete permission",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Permission Deleted Successfully",
	})
}
