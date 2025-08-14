package main

import (
	"fmt"
	"log"
	"neon/controllers"
	"neon/utilities"
	"os"

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

	oauthController := &controllers.OauthController{}

	app := echo.New()
	app.GET("/login/initiate", oauthController.Redirect)

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))))
}
