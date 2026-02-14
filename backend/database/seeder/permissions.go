package seeders

import (
	"gorm.io/gorm"

	"backend-commerce/models"
)

func SeedPermissions(db *gorm.DB) {
	permissions := []models.Permission{
		// Dashboard
		{Name: "dashboard-index"},
		{Name: "admin-dashboard"},
		{Name: "customer-dashboard"},

		// Users
		{Name: "users-index"},
		{Name: "users-create"},
		{Name: "users-show"},
		{Name: "users-edit"},
		{Name: "users-update"},
		{Name: "users-delete"},

		// Roles
		{Name: "roles-index"},
		{Name: "roles-create"},
		{Name: "roles-show"},
		{Name: "roles-edit"},
		{Name: "roles-update"},
		{Name: "roles-delete"},

		// Permissions
		{Name: "permissions-index"},
		{Name: "permissions-create"},
		{Name: "permissions-show"},
		{Name: "permissions-edit"},
		{Name: "permissions-update"},
		{Name: "permissions-delete"},

		// Categories
		{Name: "categories-index"},
		{Name: "categories-create"},
		{Name: "categories-show"},
		{Name: "categories-edit"},
		{Name: "categories-update"},
		{Name: "categories-delete"},

		// Products
		{Name: "products-index"},
		{Name: "products-create"},
		{Name: "products-show"},
		{Name: "products-edit"},
		{Name: "products-update"},
		{Name: "products-delete"},

		// Orders
		{Name: "orders-index"},
		{Name: "orders-show"},
		{Name: "orders-update"},

		// Reviews
		{Name: "reviews-index"},
		{Name: "reviews-delete"},

		// Customers
		{Name: "customers-index"},

		// Reports
		{Name: "reports-index"},

		// Sliders / Banners
		{Name: "sliders-index"},
		{Name: "sliders-create"},
		{Name: "sliders-delete"},
	}

	for _, p := range permissions {
		db.FirstOrCreate(&p, models.Permission{Name: p.Name})
	}
}
