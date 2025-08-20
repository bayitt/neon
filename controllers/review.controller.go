package controllers

import (
	"neon/services"
	"neon/validators"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReviewController struct {
	validator *validators.ReviewValidator
	service   *services.ReviewService
}

func (rc *ReviewController) create(context echo.Context) error {
	crDto, err := rc.validator.ValidateCreate(context)
	if err != nil {
		return err
	}

	review, createErr := rc.service.Create(crDto)
	if createErr != nil {
		return createErr
	}
	return context.JSON(http.StatusCreated, review)
}
