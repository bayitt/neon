package utilities

import (
	"fmt"
	"log"
	"neon/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var database *gorm.DB

func GetDatabaseObject() *gorm.DB {
	if database != nil {
		return database
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable Timezone=GMT",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("There was a problem connecting to the database")
	}

	database = db

	return database
}

func InitDatabase() {
	database := GetDatabaseObject()
	database.SetupJoinTable(&models.Category{}, "Reviews", &models.CategoryReview{})
	database.AutoMigrate(
		&models.Category{},
		&models.Series{},
		&models.Review{},
		&models.Subscriber{},
		&models.ReadingList{},
	)
}
