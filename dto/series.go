package dto

type CreateSeriesDto struct {
	Name        string  `json:"name" validate:"required,min=3"`
	Author      string  `json:"author" validate:"required,min=3"`
	Description *string `json:"description"`
}
