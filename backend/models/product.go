package models

import (
	"time"
)

type Product struct {
	Id          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug" gorm:"unique"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price"`
	Stock       int            `json:"stock"`
	CategoryId  uint           `json:"category_id"`
	Category    Category       `json:"category" gorm:"foreignKey:CategoryId"`
	Images      []ProductImage `json:"images" gorm:"foreignKey:ProductId"`
	Reviews     []Review       `json:"reviews" gorm:"foreignKey:ProductId"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
