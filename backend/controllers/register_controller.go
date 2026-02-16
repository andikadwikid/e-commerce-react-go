package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend-commerce/database"
	"backend-commerce/helpers"
	"backend-commerce/models"
	"backend-commerce/structs"

)

// Register adalah handler untuk proses registrasi user baru.
// Flow besarnya:
// 1) Ambil dan validasi data dari body request
// 2) Cek apakah email sudah pernah dipakai user lain
// 3) Hash (enkripsi satu arah) password user
// 4) Siapkan data user dan cari role "user" di database
// 5) Simpan user baru dan hubungkan dengan role "user"
// 6) Kirim response sukses berisi data user yang sudah tersimpan
func Register(c *gin.Context) {
	// LANGKAH 1: Ambil data dari body request dan validasi

	// Struct ini menjadi tempat menampung data JSON yang dikirim client
	var req structs.RegisterRequest

	// Bind JSON ke struct req sekaligus melakukan validasi sesuai tag di struct
	err := c.ShouldBindJSON(&req)
	if err != nil {
		// Jika validasi gagal, langsung balas dengan status 422 dan detail error
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err, req),
		})
		// Proses dihentikan karena data dari client tidak valid
		return
	}

	// LANGKAH 2: Cek apakah email sudah pernah didaftarkan

	// existingUser akan diisi jika email sudah ada di database
	var existingUser models.User

	// Cari user berdasarkan email, jika tidak ada error berarti user dengan email itu sudah ada
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		// Email sudah dipakai, balas dengan status 400 dan info bahwa email sudah terdaftar
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Email already registered",
			Errors:  map[string]string{"email": "Email is already in use"},
		})
		// Proses dihentikan karena email harus unik
		return
	}

	// LANGKAH 3: Hash password sebelum disimpan

	// HashPassword akan mengubah password biasa menjadi password yang aman (tidak bisa dibaca langsung)
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		// Jika proses hash gagal, berarti ada masalah di server
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to hash password",
		})
		// Proses dihentikan karena tidak boleh menyimpan password tanpa hash
		return
	}

	// LANGKAH 4: Siapkan data user baru yang akan disimpan

	user := models.User{
		// Nama dan email langsung diambil dari data request
		Name:  req.Name,
		Email: req.Email,
		// Password yang disimpan adalah yang sudah di-hash
		Password: hashedPassword,
		// Username dibuat otomatis dari nama dalam bentuk slug (huruf kecil, spasi jadi tanda hubung)
		Username: helpers.Slugify(req.Name),
	}

	// LANGKAH 5: Ambil role "user" lalu simpan user baru

	// Role akan dihubungkan ke user sebagai hak akses dasar
	var role models.Role

	// Cari role dengan nama "user"
	if err := database.DB.Where("name = ?", "user").First(&role).Error; err != nil {
		// Jika role "user" tidak ditemukan, proses tidak bisa dilanjutkan
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Role 'user' not found",
		})
		return
	}

	// Simpan user baru ke database
	if err := database.DB.Create(&user).Error; err != nil {
		// Jika gagal menyimpan user, balas dengan status 500
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
		})
		return
	}

	// Hubungkan user baru dengan role "user"
	database.DB.Model(&user).Association("Roles").Append(&role)

	// LANGKAH 6: Kirim response sukses ke client

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Registration successful",
		Data:    user,
	})
}
