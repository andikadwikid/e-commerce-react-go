package seeders

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"backend-commerce/models"
)

func SeedUsers(db *gorm.DB) {
	// Hash password default "password"
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	// Dapatkan role yang sudah di-seed sebelumnya
	var adminRole models.Role

	// Pastikan nama role sesuai dengan yang ada di SeedRoles (lowercase)
	db.Where("name = ?", "admin").First(&adminRole)

	// Data user awal (Hanya Admin)
	users := []models.User{
		{
			Name:     "Admin Toko",
			Username: "admin",
			Email:    "admin@toko.com",
			Password: string(password),
			Roles:    []models.Role{adminRole},
		},
	}

	for _, u := range users {
		var user models.User
		// Cek apakah user sudah ada berdasarkan username
		if err := db.Where("username = ?", u.Username).First(&user).Error; err != nil {
			// Jika belum ada, buat baru
			db.Create(&u)
		} else {
			// Jika sudah ada, update info dasar (password reset ke default jika seed dijalankan ulang)
			db.Model(&user).Updates(models.User{
				Name:     u.Name,
				Email:    u.Email,
				Password: string(password),
			})
			// Update juga relasi Role-nya
			db.Model(&user).Association("Roles").Replace(u.Roles)
		}
	}
}
