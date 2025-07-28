package helpers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// TranslateErrorMessage handles validation errors from validator.v10 and duplicate entry errors from GORM
func TranslateErrorMessage(err error) map[string]string {
	errorsMap := make(map[string]string)
	if err == nil {
			return errorsMap
	}

	// Handle validator.v10 errors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
					field := fieldError.Field()
					errorsMap[field] = validatorErrorMessage(fieldError)
			}
	}

	// Handle GORM duplicate entry and not found errors
	for k, v := range handleGormError(err) {
		errorsMap[k] = v
	}

	return errorsMap
}

func validatorErrorMessage(fieldError validator.FieldError) string {
    switch fieldError.Tag() {
    case "required":
        return fmt.Sprintf("%s is required", fieldError.Field())
    case "email":
        return "Invalid email format"
    case "unique":
        return fmt.Sprintf("%s already exists", fieldError.Field())
    case "min":
        return fmt.Sprintf("%s must be at least %s characters", fieldError.Field(), fieldError.Param())
    case "max":
        return fmt.Sprintf("%s must be at most %s characters", fieldError.Field(), fieldError.Param())
    case "numeric":
        return fmt.Sprintf("%s must be a number", fieldError.Field())
    default:
        return "Invalid value"
    }
}

func handleGormError(err error) map[string]string {
	errorsMap := make(map[string]string)
	if err == nil {
		return errorsMap
	}
	if strings.Contains(err.Error(), "Duplicate entry") {
		for k, v := range checkDuplicateFields(err, "username", "email") {
			errorsMap[k] = v
		}
	} else if err == gorm.ErrRecordNotFound {
		errorsMap["Error"] = "Record not found"
	}
	return errorsMap
}

func checkDuplicateFields(err error, fields ...string) map[string]string {
	errorsMap := make(map[string]string)
	if err == nil {
		return errorsMap
	}
	for _, field := range fields {
		if strings.Contains(err.Error(), field) {
			titleCaser := cases.Title(language.Und)
			errorsMap[titleCaser.String(field)] = fmt.Sprintf("%s already exists", titleCaser.String(field))
		}
	}
	return errorsMap
}

// IsDuplicateEntryError mendeteksi apakah error dari database adalah duplicate entry
func IsDuplicateEntryError(err error) bool {
	// Mengecek apakah error merupakan duplikasi entri
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}
