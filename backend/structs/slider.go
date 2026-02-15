package structs

import "mime/multipart"

type SliderCreateRequest struct {
	Link  string                `form:"link"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

type SliderResponse struct {
	Id    uint   `json:"id"`
	Link  string `json:"link"`
	Image string `json:"image"`
}
