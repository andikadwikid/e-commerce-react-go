package models

import "time"

type Address struct {
	Id            uint      `json:"id" gorm:"primaryKey"`
	UserId        uint      `json:"user_id"`
	User          User      `json:"user" gorm:"foreignKey:UserId"`
	RecipientName string    `json:"recipient_name"`
	Phone         string    `json:"phone"`
	AddressLine1  string    `json:"address_line1"`
	AddressLine2  string    `json:"address_line2"`
	District      string    `json:"district"`
	DistrictId    string    `json:"district_id"`
	City          string    `json:"city"`
	CityId        string    `json:"city_id"`
	Province      string    `json:"province"`
	ProvinceId    string    `json:"province_id"`
	PostalCode    string    `json:"postal_code"`
	IsPrimary     bool      `json:"is_primary" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
