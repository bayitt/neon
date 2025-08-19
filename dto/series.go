package dto

import "github.com/google/uuid"

type CreateSeriesDto struct {
	CategoryUuid uuid.UUID `json:"category_uuid" validate:"required"`
	Name         string    `json:"name" validate:"required,min=3"`
	Author       string    `json:"author" validate:"required,min=3"`
	Description  *string   `json:"description"`

	CategoryID uint
}
