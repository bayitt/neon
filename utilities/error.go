package utilities

import "github.com/labstack/echo/v4"

func ThrowError(code int, error string, message string) error {
	return echo.NewHTTPError(code, map[string]string{"error": error, "message": message})
}
