package services

import (
	"errors"
	"fmt"
	"neon/models"
	"neon/utilities"
	"net/http"

	"gorm.io/gorm"
)

type ReviewService struct {
	DB *gorm.DB
}

func (rs *ReviewService) FindUnique(field string, value string) (models.Review, error) {
	var review models.Review
	result := rs.DB.Where(fmt.Sprintf("%s = ?", field), value).First(&review)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Review{}, utilities.ThrowError(
			http.StatusNotFound,
			"REVIEW_002",
			fmt.Sprintf("review with field %s and value %s does not exist", field, value),
		)
	}
	return review, nil
}
