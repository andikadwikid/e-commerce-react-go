package seeders

import (
	"log"

	"backend-commerce/database"
)

// Seed menjalankan semua seeder yang terdaftar
func Seed() {
	db := database.DB
	if db == nil {
		log.Println("Database connection is nil, skipping seeding.")
		return
	}

	log.Println("Running database seeders...")

	// 1. Permissions (Harus duluan)
	log.Println("Seeding permissions...")
	SeedPermissions(db)

	// 2. Roles (Butuh Permissions)
	log.Println("Seeding roles...")
	SeedRoles(db)

	// 3. Users (Butuh Roles)
	log.Println("Seeding users...")
	SeedUsers(db)

	log.Println("Database seeding completed!")
}
