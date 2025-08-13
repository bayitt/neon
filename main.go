package main

import (
	"fmt"
	"log"
	"neon/models"
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

	database := utilities.GetDatabaseObject()
	database.AutoMigrate(&models.Category{}, &models.Series{}, &models.Review{}, &models.Subscriber{})

	app := echo.New()

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))))
}
