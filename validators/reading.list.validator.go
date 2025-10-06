package validators

import (
	"neon/dto"
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
	_, readingListErr := rlv.Service.FindUnique("title", crlDto.Title)

	if readingListErr != nil {
		return nil, readingListErr
	}

	return crlDto, nil
}
