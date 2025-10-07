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
	result := rls.DB.Where(fmt.Sprintf("%s = ?", field), value).First(&readingListItem)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return readingListItem, utilities.ThrowError(
			http.StatusNotFound,
			"READING_LIST_002",
			fmt.Sprintf(
				"book with field %s and value of %s is not present in the reading list",
				field,
				value,
			),
		)
	}

	return readingListItem, nil
}

func (rls *ReadingListService) Update(
	readingListItem models.ReadingList,
	urlDto *dto.UpdateReadingListDto,
) (models.ReadingList, error) {
	err := rls.DB.Transaction(func(tx *gorm.DB) error {
		updateStringField := func(initialValue string, param *string) string {
			if param != nil {
				return *param
			}
			return initialValue
		}

		updateBoolField := func(initialValue uint, param *bool) uint {
			if param != nil {
				if *param {
					return 1
				} else {
					return 0
				}
			}
			return initialValue
		}

		readingListItem.Title = updateStringField(readingListItem.Title, urlDto.Title)
		readingListItem.Author = updateStringField(readingListItem.Author, urlDto.Author)
		readingListItem.Image = updateStringField(readingListItem.Image, urlDto.Image)
		readingListItem.Status = updateBoolField(readingListItem.Status, urlDto.Status)

		result := rls.DB.Save(&readingListItem)

		if result.Error != nil {
			utilities.ThrowError(
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				result.Error.Error(),
			)
		}

		if readingListItem.Status != 1 {
			return nil
		}

		result = rls.DB.Model(&models.ReadingList{}).
			Where("id <> ?", readingListItem.ID).
			Update("status", 0)

		if result.Error != nil {
			utilities.ThrowError(
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				result.Error.Error(),
			)
		}

		return nil
	})

	if err != nil {
		return models.ReadingList{}, err
	}

	return readingListItem, nil
}

func (rls *ReadingListService) Create(
	crlDto *dto.CreateReadingListDto,
) (models.ReadingList, error) {
	readingListItem := models.ReadingList{
		Uuid:   crlDto.Uuid,
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
