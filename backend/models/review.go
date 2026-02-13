package models

import "time"

type Review struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	UserId    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserId"`
	ProductId uint      `json:"product_id"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductId"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
