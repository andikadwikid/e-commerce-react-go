package structs

import "time"

type ReviewCreateRequest struct {
	ProductId uint   `json:"product_id" form:"product_id" binding:"required"`
	Rating    int    `json:"rating" form:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" form:"comment" binding:"required"`
}

type ReviewResponse struct {
	Id        uint         `json:"id"`
	User      UserResponse `json:"user"`
	Product   string       `json:"product_name,omitempty"`
	Rating    int          `json:"rating"`
	Comment   string       `json:"comment"`
	CreatedAt time.Time    `json:"created_at"`
}
