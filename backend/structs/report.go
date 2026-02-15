package structs

type ReportRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type ReportResponse struct {
	TotalRevenue int64 `json:"total_revenue"`
	TotalOrders  int64 `json:"total_orders"`
}
