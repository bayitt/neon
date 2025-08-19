package main

import (
	"fmt"
	"log"
	"neon/controllers"
	"neon/utilities"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("There was a problem loading the environment variables")
	}

	utilities.InitDatabase()
	utilities.RegisterOauthProviders()

	app := echo.New()
	app.Validator = &utilities.RequestValidator{Validator: validator.New()}
	controllers.RegisterOauthRoutes(app.Group("/login"))
	controllers.RegisterCategoryRoutes(app.Group("/categories"))
	controllers.RegisterSeriesRoutes(app.Group("/series"))

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))))
}
