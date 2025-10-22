package services

import (
	"neon/dto"
	"neon/models"
	"neon/utilities"
	"net/http"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SubscriberService struct {
	DB *gorm.DB
}

func (ss *SubscriberService) Create(csDto *dto.CreateSubscriberDto) (models.Subscriber, error) {
	subscriber := models.Subscriber{Email: csDto.Email}
	result := ss.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&subscriber)

	if result.Error != nil {
		return models.Subscriber{}, utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			result.Error.Error(),
		)
	}

	return subscriber, nil
}

func (ss *SubscriberService) Find() ([]models.Subscriber, error) {
	var subscribers []models.Subscriber
	result := ss.DB.Find(&subscribers)

	if result.Error != nil {
		return []models.Subscriber{}, utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			result.Error.Error(),
		)
	}

	return subscribers, nil
}
