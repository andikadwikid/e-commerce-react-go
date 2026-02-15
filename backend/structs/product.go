package structs

import "backend-commerce/models"

type ProductCreateRequest struct {
	Name        string  `json:"name" form:"name" binding:"required"`
	Description string  `json:"description" form:"description" binding:"required"`
	Price       float64 `json:"price" form:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" form:"stock" binding:"required,gte=0"`
	Category    uint    `json:"category" form:"category" binding:"required"`
}

type ProductUpdateRequest struct {
	Name        string  `json:"name" form:"name" binding:"required"`
	Description string  `json:"description" form:"description" binding:"required"`
	Price       float64 `json:"price" form:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" form:"stock" binding:"required,gte=0"`
	Category    uint    `json:"category" form:"category" binding:"required"`
}

type ProductResponse struct {
	Id          uint                   `json:"id"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description"`
	Price       float64                `json:"price"`
	Stock       int                    `json:"stock"`
	Category    CategoryResponse       `json:"category"`
	CategoryId  uint                   `json:"category_id"`
	Images      []ProductImageResponse `json:"images"`
	CanReview   bool                   `json:"can_review"`
}

type ProductImageResponse struct {
	Id       uint   `json:"id"`
	ImageUrl string `json:"image"`
}

func ToProductResponse(product models.Product) ProductResponse {
	images := []ProductImageResponse{}
	for _, img := range product.Images {
		url := img.ImageUrl
		if len(url) > 0 && url[0] != '/' && len(url) < 4 || url[:4] != "http" {
			url = "/uplaod/products/" + url
		}
		images = append(images, ProductImageResponse{
			Id:       img.Id,
			ImageUrl: url,
		})
	}

	return ProductResponse{
		Id:          product.Id,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    ToCategoryResponse(product.Category),
		CategoryId:  product.CategoryId,
		Images:      images,
	}
}

func ToProductResponseWithBaseURL(product models.Product, baseURL string) ProductResponse {
	images := []ProductImageResponse{}
	for _, img := range product.Images {
		url := img.ImageUrl
		if len(url) > 0 {
			if url[0] != '/' && (len(url) < 4 || url[:4] != "http") {
				url = baseURL + "/uploads/products/" + url
			} else if url[0] == '/' {
				url = baseURL + url
			}
		}
		images = append(images, ProductImageResponse{
			Id:       img.Id,
			ImageUrl: url,
		})
	}

	return ProductResponse{
		Id:          product.Id,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    ToCategoryResponse(product.Category),
		CategoryId:  product.CategoryId,
		Images:      images,
	}
}
