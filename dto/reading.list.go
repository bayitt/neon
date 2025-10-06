package dto

import "github.com/google/uuid"

type CreateReadingListDto struct {
	Title  string `form:"title" validate:"required"`
	Author string `form:"title" validate:"required"`

	Uuid  uuid.UUID
	Image string
}
