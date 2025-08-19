package services

import (
	"errors"
	"fmt"
	"neon/dto"
	"neon/models"
	"strings"

	"gorm.io/gorm"
)

type SeriesService struct {
	DB *gorm.DB
}

func (ss *SeriesService) FindUnique(field string, value string) (models.Series, error) {
	var series models.Series
	result := ss.DB.Where(fmt.Sprintf("%s - ?", field), value).First(&series)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return series, fmt.Errorf("series with field %s and value %s does not exist", field, value)
	}

	return series, nil
}

func (ss *SeriesService) Create(csDto *dto.CreateSeriesDto) (*models.Series, error) {
	slug := "/" + strings.Replace(csDto.Name, " ", "/", -1)
	series := &models.Series{CategoryID: csDto.CategoryID, Name: csDto.Name, Slug: slug, Author: csDto.Author, Description: csDto.Description}
	result := ss.DB.Save(series)

	if result.Error != nil {
		return nil, errors.New(result.Error.Error())
	}

	return series, nil
}
