package seeders

import (
	"gorm.io/gorm"

	"backend-commerce/models"
)

func SeedRoles(db *gorm.DB) {
	// Data role dasar dengan nama
	roles := []models.Role{
		{Name: "admin"},
		{Name: "user"},
	}

	// Loop user assignments
	for _, role := range roles {
		// Cek/insert role
		db.FirstOrCreate(&role, models.Role{Name: role.Name})

		// Khusus Admin: Assign semua permission KECUALI yang dimulai dengan "customer-"
		if role.Name == "admin" {
			var adminPermissions []models.Permission
			db.Where("name NOT LIKE ?", "customer-%").Find(&adminPermissions)
			db.Model(&role).Association("Permissions").Replace(adminPermissions)
		}

		// Khusus User (Customer): Assign semua permission yang dimulai dengan "customer-"
		if role.Name == "user" {
			var customerPermissions []models.Permission
			db.Where("name LIKE ?", "customer-%").Find(&customerPermissions)
			db.Model(&role).Association("Permissions").Replace(customerPermissions)
		}
	}
}
