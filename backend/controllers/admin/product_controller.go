package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend-commerce/database"
	"backend-commerce/helpers"
	"backend-commerce/models"
	"backend-commerce/structs"
)

func FindProducts(c *gin.Context) {
	var products []models.Product
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)
	hostURL := helpers.BuildHostURL(c)

	query := database.DB.Model(&models.Product{}).Preload("Category").Preload("Images")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to count products",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch products",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	// Transform ke Response Struct dengan full image URL
	var data []structs.ProductResponse
	for _, p := range products {
		data = append(data, structs.ToProductResponseWithBaseURL(p, hostURL))
	}

	helpers.PaginateResponse(c, data, total, page, limit, baseURL, search, "List Data Product")
}
