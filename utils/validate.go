package utils

import "github.com/go-playground/validator/v10"

type (
	validationErrorResponse struct {
		Field string `json:"field"`
		Tag   string `json:"tag"`
		Value string `json:"value"`
	}
)

func FormatValidationError(err error) []validationErrorResponse {
	var errors []validationErrorResponse

	for _, err := range err.(validator.ValidationErrors) {
		element := &validationErrorResponse{
			Field: err.StructNamespace(),
			Tag:   err.Tag(),
			Value: err.Param(),
		}

		errors = append(errors, *element)
	}

	return errors

}
