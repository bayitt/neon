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
	Status        uint       `form:"status"`

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
	Status        *uint      `             form:"status"`

	Image      *string
	Categories *[]models.Category
	Series     *models.Series
}
