package services

import (
	"errors"
	"fmt"
	"neon/dto"
	"neon/models"
	"neon/utilities"
	"net/http"
	"strings"

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

func (rs *ReviewService) Create(crDto *dto.CreateReviewDto) (models.Review, error) {
	var review models.Review
	err := rs.DB.Transaction(func(tx *gorm.DB) error {
		slug := "/" + strings.Replace(
			crDto.Title,
			" ",
			"-",
			-1,
		) + "-" + utilities.GenerateRandomString(
			4,
		)
		review = models.Review{
			Title:   crDto.Title,
			Slug:    slug,
			Author:  crDto.Author,
			Content: crDto.Content,
			Status:  crDto.Status,
		}

		if crDto.Series != nil {
			review.SeriesID = &crDto.Series.ID
		}
		result := tx.Create(&review)

		if result.Error != nil {
			return utilities.ThrowError(
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				result.Error.Error(),
			)
		}

		associateError := tx.Model(&review).Association("Categories").Append(crDto.Categories)
		if associateError != nil {
			return utilities.ThrowError(
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				associateError.Error(),
			)
		}

		return nil
	})

	if err != nil {
		return models.Review{}, err
	}

	var categories = []*models.Category{}
	for i := 0; i < len(crDto.Categories); i++ {
		categories = append(categories, &(crDto.Categories[i]))
	}
	review.Categories = categories
	review.Series = crDto.Series

	return review, nil
}
