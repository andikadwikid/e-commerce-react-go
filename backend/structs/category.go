package structs

import "backend-commerce/models"

type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func ToCategoryResponse(category models.Category) CategoryResponse {
	return CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
		Slug: category.Slug,
	}
}
