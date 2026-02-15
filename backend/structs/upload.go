package structs

import "mime/multipart"

type UploadConfig struct {
	File           *multipart.FileHeader
	AllowedTypes   []string
	MaxSize        int64
	DestinationDir string
}

type UploadResult struct {
	FileName string
	FilePath string
	Response *ErrorResponse
}
