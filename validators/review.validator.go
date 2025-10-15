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
	Rs *services.ReviewService
	Cs *services.CategoryService
	Ss *services.SeriesService
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
	categoryUuids := strings.Split(crDto.CategoryUuids, ",")
	for i := 0; i < len(categoryUuids); i++ {
		category, categoryErr := rv.Cs.FindUnique("uuid", categoryUuids[i])
		if categoryErr != nil {
			return nil, categoryErr
		}
		categories = append(categories, category)
	}
	crDto.Categories = categories

	if crDto.SeriesUuid != nil {
		series, seriesErr := rv.Ss.FindUnique("uuid", (*crDto.SeriesUuid).String())
		if seriesErr != nil {
			return nil, seriesErr
		}
		crDto.Series = &series
	}

	crDto.Title = strings.ToLower(crDto.Title)
	crDto.Author = strings.ToLower(crDto.Author)
	_, reviewErr := rv.Rs.FindUnique("title", crDto.Title, false)

	if reviewErr != nil {
		return crDto, nil
	}

	return nil, utilities.ThrowError(
		http.StatusBadRequest,
		"REVIEW_001",
		fmt.Sprintf("review with the title %s exists already", crDto.Title),
	)
}

func (rv *ReviewValidator) ValidateUpdate(
	context echo.Context,
) (models.Review, *dto.UpdateReviewDto, error) {
	urDto := new(dto.UpdateReviewDto)
	if err := context.Bind(urDto); err != nil {
		return models.Review{}, nil, utilities.ThrowError(
			http.StatusBadRequest,
			"MALFORMED_REQUEST",
			err.Error(),
		)
	}

	if err := context.Validate(urDto); err != nil {
		return models.Review{}, nil, err
	}

	review, reviewErr := rv.Rs.FindUnique("uuid", urDto.Uuid.String(), false)
	if reviewErr != nil {
		return models.Review{}, nil, reviewErr
	}

	if urDto.Title != nil {
		parsedTitle := strings.ToLower(*urDto.Title)
		urDto.Title = &parsedTitle

		namedReview, reviewErr := rv.Rs.FindUnique("title", *urDto.Title, false)
		if reviewErr == nil && namedReview.Uuid.String() != review.Uuid.String() {
			return models.Review{}, nil, utilities.ThrowError(
				http.StatusBadRequest,
				"REVIEW_001",
				fmt.Sprintf("review with title %s already exists", *urDto.Title),
			)
		}
	}

	if urDto.Author != nil {
		parsedAuthor := strings.ToLower(*urDto.Author)
		urDto.Author = &parsedAuthor
	}

	if urDto.CategoryUuids != nil {
		categoryUuids := strings.Split(*urDto.CategoryUuids, ",")
		var categories = []models.Category{}

		for i := 0; i < len(categoryUuids); i++ {
			category, categoryErr := rv.Cs.FindUnique("uuid", categoryUuids[i])
			if categoryErr != nil {
				return models.Review{}, nil, categoryErr
			}
			categories = append(categories, category)
			urDto.Categories = &categories
		}
	}

	if urDto.SeriesUuid != nil {
		series, seriesErr := rv.Ss.FindUnique("uuid", (*urDto.SeriesUuid).String())
		if seriesErr != nil {
			return models.Review{}, nil, seriesErr
		}
		urDto.Series = &series
	}

	return review, urDto, nil
}

func (rv *ReviewValidator) ValidateGet(context echo.Context) (models.Review, error) {
	grDto := new(dto.GetReviewDto)
	if err := context.Bind(grDto); err != nil {
		return models.Review{}, utilities.ThrowError(
			http.StatusBadRequest,
			"MALFORMED_REQUEST",
			err.Error(),
		)
	}

	grDto.Slug = "/" + grDto.Slug

	review, err := rv.Rs.FindUnique("slug", grDto.Slug, true)
	if err != nil {
		return models.Review{}, err
	}

	return review, nil
}

func (rv *ReviewValidator) ValidateGetMultiple(context echo.Context) (*dto.GetReviewsDto, error) {
	grDto := new(dto.GetReviewsDto)
	if err := context.Bind(grDto); err != nil {
		return nil, utilities.ThrowError(http.StatusBadRequest, "MALFORMED_REQUEST", err.Error())
	}

	if grDto.Page == 0 {
		grDto.Page = 1
	}

	if grDto.Count == 0 {
		grDto.Count = 10
	}

	return grDto, nil
}

func (rv *ReviewValidator) ValidateGetByCategory(
	context echo.Context,
) (*dto.GetReviewsByCategoryDto, error) {
	grbcDto := new(dto.GetReviewsByCategoryDto)
	if err := context.Bind(grbcDto); err != nil {
		return nil, utilities.ThrowError(http.StatusBadRequest, "MALFORMED_REQUEST", err.Error())
	}

	if grbcDto.Page == 0 {
		grbcDto.Page = 1
	}

	if grbcDto.Count == 0 {
		grbcDto.Count = 10
	}

	category, categoryErr := rv.Cs.FindUnique("uuid", grbcDto.CategoryUuid.String())
	if categoryErr != nil {
		return nil, categoryErr
	}

	grbcDto.Category = category
	return grbcDto, nil
}

func (rv *ReviewValidator) ValidateGetByCategories(
	context echo.Context,
) (*dto.GetReviewsByCategoriesDto, error) {
	grbcDto := new(dto.GetReviewsByCategoriesDto)

	if err := context.Bind(grbcDto); err != nil {
		return nil, utilities.ThrowError(http.StatusBadRequest, "MALFORMED_REQUEST", err.Error())
	}

	if err := context.Validate(grbcDto); err != nil {
		return nil, err
	}

	categoryUuids := strings.Split(grbcDto.CategoryUuids, ".")

	var categories = []models.Category{}
	for _, categoryUuid := range categoryUuids {
		category, categoryErr := rv.Cs.FindUnique("uuid", categoryUuid)
		if categoryErr != nil {
			return nil, categoryErr
		}

		categories = append(categories, category)
	}

	grbcDto.Categories = categories

	return grbcDto, nil
}

func (rv *ReviewValidator) ValidateGetBySeries(
	context echo.Context,
) (*dto.GetReviewsBySeriesDto, error) {
	grbsDto := new(dto.GetReviewsBySeriesDto)
	if err := context.Bind(grbsDto); err != nil {
		return nil, utilities.ThrowError(http.StatusBadRequest, "MALFORMED_REQUEST", err.Error())
	}

	series, seriesErr := rv.Ss.FindUnique("uuid", grbsDto.SeriesUuid.String())
	if seriesErr != nil {
		return nil, seriesErr
	}

	grbsDto.Series = series
	return grbsDto, nil
}
