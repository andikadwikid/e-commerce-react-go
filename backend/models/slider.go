package models

import "time"

type Slider struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	Image     string    `gorm:"type:varchar(255)" json:"image"`
	Link      string    `gorm:"type:varchar(255)" json:"link"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
