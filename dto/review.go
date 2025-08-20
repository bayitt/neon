package dto

import "github.com/google/uuid"

type CreateReviewDto struct {
	CategoryUuid uuid.UUID  `form:"category_uuid" validate:"required"`
	SeriesUuid   *uuid.UUID `form:"series_uuid"`
	Title        string     `form:"title"         validate:"required"`
	Content      string     `form:"content"       validate:"required"`
	Status       uint       `form:"status"`

	CategoryID uint
	SeriesID   *uint
}
