package dto

import "github.com/google/uuid"

type CreateReadingListDto struct {
	Title  string `form:"title"  validate:"required"`
	Author string `form:"author" validate:"required"`

	Uuid  uuid.UUID
	Image string
}

type UpdateReadingListDto struct {
	Uuid   uuid.UUID `param:"uuid"`
	Title  *string   `             form:"title"`
	Author *string   `             form:"author"`
	Status *bool     `             form:"status"`

	Image *string
}
