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

type SeriesService struct {
	DB *gorm.DB
}

func (ss *SeriesService) FindUnique(field string, value string) (models.Series, error) {
	var series models.Series
	result := ss.DB.Where(fmt.Sprintf("%s = ?", field), value).First(&series)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return series, utilities.ThrowError(
			http.StatusNotFound,
			"SERIES_002",
			fmt.Sprintf("series with field %s and value %s does not exist", field, value),
		)
	}

	return series, nil
}

func (ss *SeriesService) Create(csDto *dto.CreateSeriesDto) (models.Series, error) {
	slug := "/" + strings.Replace(
		csDto.Name,
		" ",
		"-",
		-1,
	) + "-" + utilities.GenerateRandomString(
		4,
	)
	series := models.Series{
		Name:        csDto.Name,
		Slug:        slug,
		Author:      csDto.Author,
		Description: csDto.Description,
	}
	result := ss.DB.Save(&series)

	if result.Error != nil {
		return series, utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			result.Error.Error(),
		)
	}

	return series, nil
}

func (ss *SeriesService) Update(
	series models.Series,
	usDto *dto.UpdateSeriesDto,
) (models.Series, error) {
	updateValue := func(initialValue string, newValuePointer *string) string {
		if newValuePointer != nil && len(*newValuePointer) > 0 {
			return strings.ToLower(*newValuePointer)
		}
		return initialValue
	}

	updateDescription := func(initialDescriptionPointer *string, newDescriptionPointer *string) *string {
		if newDescriptionPointer != nil && len(*newDescriptionPointer) > 0 {
			return newDescriptionPointer
		}
		return initialDescriptionPointer
	}

	if usDto.Name != nil && series.Name != strings.ToLower(*usDto.Name) {
		series.Slug = "/" + strings.Replace(
			strings.ToLower(*usDto.Name),
			" ",
			"-",
			-1,
		) + "-" + utilities.GenerateRandomString(
			4,
		)
	}
	series.Name = updateValue(series.Name, usDto.Name)
	series.Author = updateValue(series.Author, usDto.Author)
	series.Description = updateDescription(series.Description, usDto.Description)
	result := ss.DB.Save(&series)

	if result.Error != nil {
		return series, utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			result.Error.Error(),
		)
	}

	return series, nil
}

func (ss *SeriesService) Find(offset uint, count uint) ([]models.Series, error) {
	var series []models.Series
	result := ss.DB.Order("created_at desc").
		Offset(int(offset)).
		Limit(int(count)).
		Preload("Reviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("reviews.created_at DESC").Select("SeriesID", "Image")
		}).
		Find(&series)

	if result.Error != nil {
		return []models.Series{}, utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			result.Error.Error(),
		)
	}

	var parsedSeries []models.Series

	for _, series := range series {
		var images []string

		if len(series.Reviews) == 0 {
			series.Images = make([]string, 0)
			parsedSeries = append(parsedSeries, series)
			continue
		}

		for _, review := range series.Reviews {
			if len(*review.Image) > 0 {
				images = append(images, *review.Image)
			}
		}

		series.Images = images
		parsedSeries = append(parsedSeries, series)
	}

	return parsedSeries, nil
}

func (ss *SeriesService) Count() uint {
	var totalSeries int64
	ss.DB.Model(models.Series{}).Count(&totalSeries)

	return uint(totalSeries)
}
