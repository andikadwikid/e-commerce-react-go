package models

import "time"

type Cart struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	UserId    uint      `json:"user_id" gorm:"index"`
	User      User      `json:"user" gorm:"foreignKey:UserId"`
	ProductId uint      `json:"product_id"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductId"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
