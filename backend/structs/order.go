package structs

type CheckoutRequest struct {
	ShippingName    string  `json:"shipping_name" binding:"required"`
	ShippingPhone   string  `json:"shipping_phone" binding:"required"`
	ShippingAddress string  `json:"shipping_address" binding:"required"`
	ShippingCost    float64 `json:"shipping_cost" binding:"required"`
	Courier         string  `json:"courier" binding:"required"`
	Service         string  `json:"service" binding:"required"`
}

type CheckoutResponse struct {
	SnapToken   string  `json:"snap_token"`
	RedirectURL string  `json:"redirect_url"`
	OrderID     string  `json:"order_id"`
	TotalPrice  float64 `json:"total_price"`
}
