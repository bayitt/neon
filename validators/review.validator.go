package validators

import (
	"neon/dto"
	"neon/services"
	"neon/utilities"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ReviewValidator struct {
	rs *services.ReviewService
	cs *services.CategoryService
	ss *services.SeriesService
}

func (rv *ReviewValidator) ValidateCreate(context echo.Context) (*dto.CreateReviewDto, error) {
	crDto := new(dto.CreateReviewDto)
	if err := context.Bind(crDto); err != nil {
		return nil, utilities.ThrowError(http.StatusBadRequest, "REVIEW_001", err.Error())
	}

	if err := context.Validate(crDto); err != nil {
		return nil, err
	}

	category, categoryErr := rv.cs.FindUnique("uuid", crDto.CategoryUuid.String())
	if categoryErr != nil {
		return nil, categoryErr
	}
	crDto.CategoryID = category.ID

	if crDto.SeriesUuid != nil {
		series, seriesErr := rv.ss.FindUnique("uuid", (*crDto.SeriesUuid).String())
		if seriesErr != nil {
			return nil, seriesErr
		}
		crDto.SeriesID = &series.ID
	}

	crDto.Title = strings.ToLower(crDto.Title)
	return nil, categoryErr
}
