package structs

import "backend-commerce/models"

type CartRequest struct {
	ProductId uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required, min=1"`
}

type CartUpdateRequest struct {
	Quantity int `json:"quantity" binding:"required, min=1"`
}

type CartResponse struct {
	Id         uint            `json:"id"`
	ProductId  ProductResponse `json:"product_id"`
	Quantity   int             `json:"quantity"`
	TotalPrice float64         `json:"total_price"`
}

func ToCartResponse(cart models.Cart) CartResponse {
	return CartResponse{
		Id:         cart.Id,
		ProductId:  ToProductResponse(cart.Product),
		Quantity:   cart.Quantity,
		TotalPrice: cart.Product.Price * float64(cart.Quantity),
	}
}
