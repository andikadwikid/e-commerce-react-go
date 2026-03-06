package admin

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"backend-commerce/database"
	"backend-commerce/helpers"
	"backend-commerce/models"
	"backend-commerce/structs"
)

// FindSliders
func FindSliders(c *gin.Context) {
	var sliders []models.Slider

	if err := database.DB.Order("id desc").Find(&sliders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch sliders",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	// Transform image URLs
	sliderResponses := []structs.SliderResponse{}
	for _, s := range sliders {
		sliderResponses = append(sliderResponses, structs.SliderResponse{
			Id:    s.Id,
			Image: fmt.Sprintf("%s/uploads/sliders/%s", helpers.BuildHostURL(c), s.Image),
			Link:  s.Link,
		})
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List Data Sliders",
		Data:    sliderResponses,
	})
}

func CreateSlider(c *gin.Context) {
	var request structs.SliderCreateRequest

	// 1. Validasi Input (Form Data)
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Failed",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	// 2. Handle Image Upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Image is required",
		})
		return
	}

	// 3. Create directory if not exists
	uploadPath := "./public/uploads/sliders"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, os.ModePerm)
	}

	// 4. Save file
	// Format nama file: timestamp-namafileasli
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	filepath := filepath.Join(uploadPath, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to upload image",
		})
		return
	}

	// 5. Simpan ke Database
	slider := models.Slider{
		Image: filename,
		Link:  request.Link,
	}

	if err := database.DB.Create(&slider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create slider",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Slider Created Successfully",
		Data:    slider,
	})
}

// DeleteSlider
func DeleteSlider(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var slider models.Slider

	if err := database.DB.First(&slider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Slider Not Found",
		})
		return
	}

	// Delete File
	os.Remove(fmt.Sprintf("./public/uploads/sliders/%s", slider.Image))

	if err := database.DB.Delete(&slider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete slider",
			Errors:  helpers.TranslateErrorMessage(err, nil),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Slider Deleted Successfully",
	})
}
