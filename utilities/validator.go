package utilities

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// RequestValidator is a custom struct used to define validators for
// the Echo struct pointer in the Chequer app
type RequestValidator struct {
	Validator *validator.Validate
}

// RequestError is a custom struct used to define request body errors
// for the Chequer app
type RequestError struct {
	Param   string
	Message string
}

func getErrorMessage(error validator.FieldError) string {
	switch error.Tag() {
	case "email":
		return "A valid email is required."
	case "required":
		return "This field is required."
	case "min":
		return "Minimum 8 characters is required."
	default:
		return error.Error()
	}
}

func (requestValidator *RequestValidator) Validate(i interface{}) error {
	if err := requestValidator.Validator.Struct(i); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		requestErrors := make([]RequestError, len(validationErrors))

		for index, error := range validationErrors {
			requestErrors[index] = RequestError{
				Param:   strings.ToLower(error.Field()),
				Message: getErrorMessage(error),
			}
		}
		return echo.NewHTTPError(
			http.StatusBadRequest,
			map[string][]RequestError{"errors": requestErrors},
		)
	}

	return nil
}
