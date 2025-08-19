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
		return series, utilities.ThrowError(http.StatusNotFound, "SERIES_002", fmt.Sprintf("series with field %s and value %s does not exist", field, value))
	}

	return series, nil
}

func (ss *SeriesService) Create(csDto *dto.CreateSeriesDto) (models.Series, error) {
	slug := "/" + strings.Replace(csDto.Name, " ", "-", -1) + "-" + utilities.GenerateRandomString(4)
	series := models.Series{Name: csDto.Name, Slug: slug, Author: csDto.Author, Description: csDto.Description}
	result := ss.DB.Save(&series)

	if result.Error != nil {
		return series, utilities.ThrowError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", result.Error.Error())
	}

	return series, nil
}

func (ss *SeriesService) Update(series models.Series, usDto *dto.UpdateSeriesDto) (models.Series, error) {
	updateValue := func(initialValue string, newValue string) string {
		if len(newValue) > 0 {
			return newValue
		}
		return initialValue
	}

	updateDescription := func(initialDescription string, newDescription string) *string {
		if len(newDescription) > 0 {
			return &newDescription
		}
		return &initialDescription
	}

	if series.Name != *usDto.Name {
		series.Slug = "/" + strings.Replace(*usDto.Name, " ", "-", -1) + "-" + utilities.GenerateRandomString(4)
	}
	series.Name = updateValue(series.Name, *usDto.Name)
	series.Author = updateValue(series.Author, *usDto.Author)
	series.Description = updateDescription(*series.Description, *usDto.Description)
	result := ss.DB.Save(&series)

	if result.Error != nil {
		return series, utilities.ThrowError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", result.Error.Error())
	}

	return series, nil
}
