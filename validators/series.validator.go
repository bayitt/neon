package validators

import (
	"fmt"
	"neon/dto"
	"neon/services"
	"neon/utilities"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type SeriesValidator struct {
	Service *services.SeriesService
}

func (sv *SeriesValidator) ValidateCreate(context echo.Context) (*dto.CreateSeriesDto, error) {
	csDto := new(dto.CreateSeriesDto)

	if err := context.Bind(csDto); err != nil {
		return nil, utilities.ThrowError(http.StatusBadRequest, "MALFORMED_REQUEST", err.Error())
	}

	if err := context.Validate(csDto); err != nil {
		return nil, err
	}

	csDto.Name = strings.ToLower(csDto.Name)
	csDto.Author = strings.ToLower(csDto.Author)
	_, err := sv.Service.FindUnique("name", csDto.Name)

	if err != nil {
		return csDto, nil
	}

	return nil, utilities.ThrowError(http.StatusBadRequest, "SERIES_001", fmt.Sprintf("series with name %s exists already", csDto.Name))
}
