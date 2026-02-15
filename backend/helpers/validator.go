package helpers

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

)

func TranslateErrorMessage(err error, request interface{}) map[string]string {
	errorsMap := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var typ reflect.Type

		val := reflect.ValueOf(request)
		if val.IsValid() && (val.Kind() != reflect.Ptr || !val.IsNil()) {
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			typ = val.Type()
		}

		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			jsonTag := cleanFieldName(fieldName)
			displayName := jsonTag

			if typ != nil {
				if field, found := typ.FieldByName(fieldName); found {
					if t := field.Tag.Get("json"); t != "" && t != "-" {
						jsonTag = strings.Split(t, ",")[0]
					}
					if l := field.Tag.Get("label"); l != "" {
						displayName = l
					}
				}
			}

			errorsMap[jsonTag] = formatErrorMessage(fieldError, displayName)
		}
	} else {
		if err != nil {
			errStr := err.Error()
			if strings.Contains(errStr, "Duplicate entry") {
				field := extractDuplicateField(errStr)
				if field != "" {
					errorsMap[field] = fmt.Sprintf("%s sudah digunakan", field)
				} else {
					errorsMap["error"] = "Data duplikat terdeteksi"
				}
			} else if err == gorm.ErrRecordNotFound {
				errorsMap["error"] = "Data tidak ditemukan"
			} else {
				errorsMap["error"] = errStr
			}
		}
	}

	return errorsMap
}

func formatErrorMessage(fieldError validator.FieldError, displayName string) string {
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", displayName)
	case "email":
		return "Format email tidak valid"
	case "unique":
		return fmt.Sprintf("%s sudah digunakan", displayName)
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", displayName, fieldError.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", displayName, fieldError.Param())
	case "numeric":
		return fmt.Sprintf("%s harus berupa angka", displayName)
	case "gt":
		return fmt.Sprintf("%s harus lebih besar dari 0", displayName)
	case "gte":
		return fmt.Sprintf("%s harus lebih besar atau sama dengan 0", displayName)
	default:
		return fmt.Sprintf("Nilai %s tidak valid", displayName)
	}
}

func IsDuplicateEntryError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}

func extractDuplicateField(errMsg string) string {
	re := regexp.MustCompile(`for key '([^']+)'`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) >= 2 {
		keyName := matches[1]

		if strings.Contains(keyName, ".") {
			parts := strings.Split(keyName, ".")
			if len(parts) > 1 {
				return cleanFieldName(parts[len(parts)-1])
			}
		}

		return cleanFieldName(keyName)
	}
	return ""
}

func cleanFieldName(field string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(field, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
