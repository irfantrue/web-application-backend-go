// Package validator provides utility functions to validate structs using go-playground/validator.
package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/response"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates the fields of a struct based on its validator tags.
// It returns a slice of InvalidParam if validation fails, or nil if successful.
func ValidateStruct(data any) []response.InvalidParam {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	// Safely assert the error type to avoid panics on InvalidValidationError.
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return []response.InvalidParam{
			{
				Name:   "payload",
				Reason: "Invalid request payload structure.",
			},
		}
	}

	var errors []response.InvalidParam
	for _, err := range validationErrs {
		errors = append(errors, response.InvalidParam{
			Name:   err.Field(),
			Reason: getErrorMessage(err),
		})
	}

	return errors
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is mandatory."
	case "email":
		return "Must be a valid email address."
	case "min":
		return fmt.Sprintf("Must be at least %s characters long.", err.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters long.", err.Param())
	}

	return fmt.Sprintf("Failed validation: %s", err.Tag())
}
