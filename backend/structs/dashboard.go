package structs

type DashboardResponse struct {
	TotalRevenue   float64 `json:"total_revenue"`
	TotalOrders    int64   `json:"total_orders"`
	TotalProducts  int64   `json:"total_products"`
	TotalCustomers int64   `json:"total_customers"`
	PendingOrders  int64   `json:"pending_orders"`
	PaidOrders     int64   `json:"paid_orders"`
}

type OrderResponse struct {
	Id        string  `json:"id"`
	Customer  string  `json:"customer"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}
