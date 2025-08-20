package dto

import (
	"neon/models"

	"github.com/google/uuid"
)

type CreateReviewDto struct {
	CategoryUuids []uuid.UUID `form:"category_uuid" validate:"required"`
	SeriesUuid    *uuid.UUID  `form:"series_uuid"`
	Title         string      `form:"title"         validate:"required"`
	Author        string      `form:"content"       validate:"required"`
	Content       string      `form:"content"       validate:"required"`
	Status        uint        `form:"status"`

	Categories []models.Category
	Series     *models.Series
}
