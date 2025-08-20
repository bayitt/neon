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

	var categories = []models.Category{}
	for i := 0; i < len(crDto.CategoryUuids); i++ {
		category, categoryErr := rv.cs.FindUnique("uuid", crDto.CategoryUuids[i].String())
		if categoryErr != nil {
			return nil, categoryErr
		}
		categories = append(categories, category)
	}
	crDto.Categories = categories

	if crDto.SeriesUuid != nil {
		series, seriesErr := rv.ss.FindUnique("uuid", (*crDto.SeriesUuid).String())
		if seriesErr != nil {
			return nil, seriesErr
		}
		crDto.Series = &series
	}

	crDto.Title = strings.ToLower(crDto.Title)
	crDto.Author = strings.ToLower(crDto.Author)
	_, reviewErr := rv.rs.FindUnique("title", crDto.Title)

	if reviewErr != nil {
		return crDto, nil
	}

	return nil, utilities.ThrowError(
		http.StatusBadRequest,
		"REVIEW_001",
		fmt.Sprintf("review with the title %s exists already", crDto.Title),
	)
}
