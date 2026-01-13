package ingredient

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
		errors = append(errors, ValidationError{
			Field:   jsonErr.Field,
			Message: fmt.Sprintf("must be %s", jsonErr.Type.String()),
		})
		return errors
	}

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors = append(errors, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Message: formatValidationTag(e),
			})
		}
		return errors
	}

	errors = append(errors, ValidationError{
		Field:   "request",
		Message: err.Error(),
	})
	return errors
}

func formatValidationTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "min":
		return fmt.Sprintf("must be at least %s", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", e.Param())
	case "email":
		return "must be a valid email"
	case "url":
		return "must be a valid URL"
	default:
		return fmt.Sprintf("failed validation: %s", e.Tag())
	}
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v
	}
}
