package models

import (
	"time"
)

type Payment struct {
	Id                uint      `json:"id" gorm:"primaryKey"`
	OrderId           string    `json:"order_id" gorm:"type:varchar(50);unique"`
	TransactionId     string    `json:"transaction_id" gorm:"type:varchar(255)"`
	PaymentType       string    `json:"payment_type" gorm:"type:varchar(50)"`
	GrossAmount       float64   `json:"gross_amount"`
	TransactionStatus string    `json:"transaction_status" gorm:"type:varchar(50)"`
	FraudStatus       string    `json:"fraud_status" gorm:"type:varchar(50)"`
	TransactionTime   time.Time `json:"transaction_time"`
	Payload           string    `json:"payload" gorm:"type:text"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
