package services

import (
	"errors"
	"fmt"
	"neon/dto"
	"neon/models"
	"neon/utilities"
	"net/http"

	"gorm.io/gorm"
)

type ReadingListService struct {
	DB *gorm.DB
}

func (rls *ReadingListService) FindUnique(field string, value string) (models.ReadingList, error) {
	var readingListItem models.ReadingList
	result := rls.DB.Where(fmt.Sprintf("%s = ", field), value).First(&readingListItem)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return readingListItem, utilities.ThrowError(
			http.StatusNotFound,
			"READING_LIST_002",
			result.Error.Error(),
		)
	}

	return readingListItem, nil
}

func (rls *ReadingListService) Create(
	crlDto *dto.CreateReadingListDto,
) (models.ReadingList, error) {
	readingListItem := models.ReadingList{
		Title:  crlDto.Title,
		Author: crlDto.Author,
		Image:  crlDto.Image,
	}
	result := rls.DB.Save(&readingListItem)

	if result.Error != nil {
		return readingListItem, utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			result.Error.Error(),
		)
	}

	return readingListItem, nil
}
