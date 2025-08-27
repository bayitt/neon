package dto

type CreateSubscriberDto struct {
	Email string `json:"email" validate:"required,email"`
}
