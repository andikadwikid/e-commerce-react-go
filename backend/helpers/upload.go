package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend-commerce/structs"

)

func UploadFile(c *gin.Context, config structs.UploadConfig) structs.UploadResult {
	if config.File == nil {
		return structs.UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "File is required",
				Errors:  map[string]string{"file": "No file was uploaded"},
			},
		}
	}

	if config.File.Size > config.MaxSize {
		return structs.UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "File too large",
				Errors:  map[string]string{"file": fmt.Sprintf("Maximum file size is: %dMB", config.MaxSize/(1<<20))},
			},
		}
	}

	ext := strings.ToLower(filepath.Ext(config.File.Filename))
	allowed := false
	for _, t := range config.AllowedTypes {
		if ext == t {
			allowed = true
			break
		}
	}
	if !allowed {
		return structs.UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Invalid file type",
				Errors:  map[string]string{"file": fmt.Sprintf("Allowed file types: %v", config.AllowedTypes)},
			},
		}
	}

	uuidName := uuid.New().String()
	filename := uuidName + ext
	filePath := filepath.Join(config.DestinationDir, filename)

	if err := os.MkdirAll(config.DestinationDir, 0755); err != nil {
		return structs.UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Failed to create upload directory",
				Errors:  map[string]string{"system": err.Error()},
			},
		}
	}

	if err := c.SaveUploadedFile(config.File, filePath); err != nil {
		return structs.UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Failed to save file",
				Errors:  map[string]string{"file": err.Error()},
			},
		}
	}

	return structs.UploadResult{
		FileName: filename,
		FilePath: filePath,
	}
}

func RemoveFile(filePath string) error {
	return os.Remove(filePath)
}
