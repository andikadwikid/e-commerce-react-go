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

func CreateProduct(c *gin.Context) {
	var request structs.ProductCreateRequest

	// Gunakan ShouldBind agar support multipart/form-data dan validasi struct
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Failed",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	// Buat Slug
	slug := helpers.Slugify(request.Name)

	product := models.Product{
		Name:        request.Name,
		Slug:        slug,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		CategoryId:  request.CategoryId,
	}

	// Save Product
	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create product",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["images[]"]

	if len(files) > 0 {
		config := structs.UploadConfig{
			AllowedTypes:   []string{".jpg", ".jpeg", ".png"},
			MaxSize:        1024 * 1024 * 2, // 2MB
			DestinationDir: "./public/uploads/products",
		}

		for _, file := range files {
			config.File = file
			res := helpers.UploadFile(c, config)

			if res.Response == nil {
				// Insert DB
				productImage := models.ProductImage{
					ProductId: product.Id,
					ImageUrl:  res.FileName,
					IsPrimary: true,
				}

				database.DB.Create(&productImage)
			}
		}
	}

	database.DB.Preload("Category").Preload("Images").First(&product, product.Id)

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Product created successfully",
		Data:    structs.ToProductResponse(product),
	})
}
