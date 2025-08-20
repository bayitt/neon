package dto

import "github.com/google/uuid"

type CreateCategoryDto struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryDto struct {
	CreateCategoryDto
	Uuid uuid.UUID `param:"uuid"`
}
