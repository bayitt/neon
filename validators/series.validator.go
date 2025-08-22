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

	return nil, utilities.ThrowError(
		http.StatusBadRequest,
		"SERIES_001",
		fmt.Sprintf("series with name %s exists already", csDto.Name),
	)
}

func (sv *SeriesValidator) ValidateUpdate(
	context echo.Context,
) (models.Series, *dto.UpdateSeriesDto, error) {
	usDto := new(dto.UpdateSeriesDto)
	if err := context.Bind(usDto); err != nil {
		return models.Series{}, nil, utilities.ThrowError(
			http.StatusBadRequest,
			"MALFORMED_REQUEST",
			err.Error(),
		)
	}

	if err := context.Validate(usDto); err != nil {
		return models.Series{}, nil, err
	}

	series, err := sv.Service.FindUnique("uuid", usDto.SeriesUuid.String())
	if err != nil {
		return models.Series{}, nil, utilities.ThrowError(
			http.StatusNotFound,
			"SERIES_002",
			fmt.Sprintf("series with uuid %s does not exist", usDto.SeriesUuid.String()),
		)
	}

	if usDto.Name != nil && len(*usDto.Name) > 0 {
		namedSeries, err := sv.Service.FindUnique("name", *usDto.Name)
		if err == nil && namedSeries.Uuid.String() != usDto.SeriesUuid.String() {
			return models.Series{}, nil, utilities.ThrowError(
				http.StatusBadRequest,
				"SERIES_001",
				fmt.Sprintf("series with name %s exists already", *usDto.Name),
			)
		}
	}

	return series, usDto, nil
}
