package structs

type OrderAdminResponse struct {
	ID          string `json:"id"`
	Invoice     string `json:"invoice"`
	Customer    string `json:"customer"`
	Total       int    `json:"total"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	PaymentType string `json:"payment_type"`
}

type OrderAdminFilter struct {
	Status string `form:"status"`
	Date   string `form:"date"`
}
