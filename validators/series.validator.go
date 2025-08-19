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
	Cs *services.CategoryService
	Ss *services.SeriesService
}

func (sv *SeriesValidator) ValidateCreate(context echo.Context) (*dto.CreateSeriesDto, error) {
	csDto := new(dto.CreateSeriesDto)

	if err := context.Bind(csDto); err != nil {
		return nil, utilities.ThrowError(http.StatusBadRequest, "MALFORMED_REQUEST", err.Error())
	}

	if err := context.Validate(csDto); err != nil {
		return nil, err
	}

	category, categoryErr := sv.Cs.FindUnique("uuid", csDto.CategoryUuid.String())
	if categoryErr != nil {
		return nil, utilities.ThrowError(http.StatusNotFound, "CATEGORY_001", fmt.Sprintf("category with uuid %s does not exist", csDto.CategoryUuid.String()))
	}

	csDto.CategoryID = category.ID
	csDto.Name = strings.ToLower(csDto.Name)
	csDto.Author = strings.ToLower(csDto.Author)
	_, err := sv.Ss.FindUnique("name", csDto.Name)

	if err != nil {
		return csDto, nil
	}

	return nil, utilities.ThrowError(http.StatusBadRequest, "SERIES_001", fmt.Sprintf("series with name %s exists already", csDto.Name))
}
