package models

import "time"

type ProductImage struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	ProductId uint      `json:"product_id"`
	ImageUrl  string    `json:"image_url"`
	IsPrimary bool      `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
