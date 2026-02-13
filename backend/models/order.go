package models

import (
	"time"
)

type Order struct {
	Id              string      `json:"id" gorm:"primaryKey;type:varchar(50)"`
	UserId          uint        `json:"user_id"`
	User            User        `json:"user" gorm:"foreignKey:UserId"`
	TotalPrice      float64     `json:"total_price"`
	Status          string      `json:"status" gorm:"type:varchar(20);default:'pending'"`
	ShippingName    string      `json:"shipping_name" gorm:"type:varchar(255)"`
	ShippingPhone   string      `json:"shipping_phone" gorm:"type:varchar(20)"`
	ShippingAddress string      `json:"shipping_address" gorm:"type:text"`
	ShippingCost    float64     `json:"shipping_cost"`
	Courier         string      `json:"courier" gorm:"type:varchar(50)"`
	Service         string      `json:"service" gorm:"type:varchar(50)"`
	SnapToken       string      `json:"snap_token" gorm:"type:varchar(255)"`
	SnapRedirectUrl string      `json:"snap_redirect_url" gorm:"type:varchar(255)"`
	Items           []OrderItem `json:"items" gorm:"foreignKey:OrderId"`
	Payment         *Payment    `json:"payment,omitempty" gorm:"foreignKey:OrderId"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	Id        uint    `json:"id" gorm:"primaryKey"`
	OrderId   string  `json:"order_id" gorm:"type:varchar(50);index"`
	ProductId uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	SubTotal  float64 `json:"sub_total"`
}
