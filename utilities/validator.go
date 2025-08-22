package utilities

import (
	"bytes"
	"net/http"
	"unicode"

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
		return "a valid email is required."
	case "required":
		return "this field is required."
	case "min":
		return "minimum 8 characters is required."
	default:
		return error.Error()
	}
}

func toSnakeCase(s string) string {
	buf := new(bytes.Buffer)
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func (requestValidator *RequestValidator) Validate(i interface{}) error {
	if err := requestValidator.Validator.Struct(i); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		requestErrors := make([]RequestError, len(validationErrors))

		for index, error := range validationErrors {
			requestErrors[index] = RequestError{
				Param:   toSnakeCase(error.Field()),
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
