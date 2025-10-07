package utilities

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterCors(app *echo.Echo) {
	corsOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOrigins,
		AllowHeaders: []string{"*"},
	}))
}
