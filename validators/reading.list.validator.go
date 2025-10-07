package validators

import (
	"fmt"
	"neon/dto"
	"neon/models"
	"neon/services"
	"neon/utilities"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ReadingListValidator struct {
	Service *services.ReadingListService
}

func (rlv *ReadingListValidator) ValidateCreate(
	context echo.Context,
) (*dto.CreateReadingListDto, error) {
	crlDto := new(dto.CreateReadingListDto)

	if err := context.Bind(crlDto); err != nil {
		return nil, utilities.ThrowError(http.StatusBadRequest, "MALFORMED_REQUEST", err.Error())
	}

	if err := context.Validate(crlDto); err != nil {
		return nil, err
	}

	crlDto.Title = strings.ToLower(crlDto.Title)
	crlDto.Author = strings.ToLower(crlDto.Author)
	_, readingListErr := rlv.Service.FindUnique("title", crlDto.Title)

	if readingListErr != nil {
		return crlDto, nil
	}

	return crlDto, utilities.ThrowError(
		http.StatusBadRequest,
		"READING_LIST_001",
		fmt.Sprintf("%s is already present in the reading list", crlDto.Title),
	)
}

func (rlv *ReadingListValidator) ValidateUpdate(
	context echo.Context,
) (models.ReadingList, *dto.UpdateReadingListDto, error) {
	urlDto := new(dto.UpdateReadingListDto)

	if err := context.Bind(urlDto); err != nil {
		return models.ReadingList{}, nil, utilities.ThrowError(
			http.StatusBadRequest,
			"MALFORMED_REQUEST",
			err.Error(),
		)
	}

	readingListItem, readingListErr := rlv.Service.FindUnique("uuid", urlDto.Uuid.String())
	if readingListErr != nil {
		return models.ReadingList{}, nil, readingListErr
	}

	if urlDto.Title != nil {
		parsedTitle := strings.ToLower(*urlDto.Title)
		urlDto.Title = &parsedTitle
	}

	if urlDto.Author != nil {
		parsedAuthor := strings.ToLower(*urlDto.Author)
		urlDto.Author = &parsedAuthor
	}

	return readingListItem, urlDto, nil
}
