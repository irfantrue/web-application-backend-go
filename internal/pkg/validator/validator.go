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

func ValidateStruct(data any) []response.InvalidParam {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	var errors []response.InvalidParam

	for _, err := range err.(validator.ValidationErrors) {
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
	}

	return fmt.Sprintf("Failed validation: %s", err.Tag())
}
