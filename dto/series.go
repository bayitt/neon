package dto

import "github.com/google/uuid"

type CreateSeriesDto struct {
	Name        string  `json:"name"        validate:"required,min=3"`
	Author      string  `json:"author"      validate:"required,min=3"`
	Description *string `json:"description"`
}

type UpdateSeriesDto struct {
	SeriesUuid  uuid.UUID `param:"uuid"`
	Name        *string   `             json:"name"`
	Author      *string   `             json:"author"`
	Description *string   `             json:"description"`
}

type GetSeriesDto struct {
	Page  uint `query:"page"`
	Count uint `query:"count"`
}
