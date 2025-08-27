package dto

import (
	"neon/models"

	"github.com/google/uuid"
)

type CreateReviewDto struct {
	CategoryUuids string     `form:"category_uuids" validate:"required"`
	SeriesUuid    *uuid.UUID `form:"series_uuid"`
	Title         string     `form:"title"          validate:"required"`
	Author        string     `form:"author"         validate:"required"`
	Content       string     `form:"content"        validate:"required"`
	Status        bool       `form:"status"`

	Uuid       uuid.UUID
	Image      *string
	Categories []models.Category
	Series     *models.Series
}

type UpdateReviewDto struct {
	Uuid          uuid.UUID  `param:"uuid"`
	CategoryUuids *string    `             form:"category_uuids"`
	SeriesUuid    *uuid.UUID `             form:"series_uuid"`
	Title         *string    `             form:"title"`
	Author        *string    `             form:"author"`
	Content       *string    `             form:"content"`
	Status        *bool      `             form:"status"`

	Image      *string
	Categories *[]models.Category
	Series     *models.Series
}

type GetReviewDto struct {
	Uuid uuid.UUID `param:"uuid"`
}

type GetReviewsDto struct {
	Page  uint `query:"page"`
	Count uint `query:"count"`
}

type GetReviewsByCategoryDto struct {
	GetReviewsDto
	CategoryUuid uuid.UUID `param:"category_uuid"`

	Category models.Category
}

type GetReviewsBySeriesDto struct {
	SeriesUuid uuid.UUID `param:"series_uuid"`

	Series models.Series
}
