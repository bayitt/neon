package controllers

import (
	"neon/services"
	"neon/utilities"
	"neon/validators"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SubscriberController struct {
	service   *services.SubscriberService
	validator *validators.SubscriberValidator
}

func RegisterSubscriberRoutes(group *echo.Group) {
	db := utilities.GetDatabaseObject()
	ss := &services.SubscriberService{DB: db}
	sv := &validators.SubscriberValidator{}
	sc := &SubscriberController{service: ss, validator: sv}

	group.POST("", sc.Create)
}

func (sc *SubscriberController) Create(context echo.Context) error {
	csDto, err := sc.validator.ValidateCreate(context)
	if err != nil {
		return err
	}

	subscriber, subscriberErr := sc.service.Create(csDto)
	if subscriberErr != nil {
		return subscriberErr
	}

	return context.JSON(http.StatusCreated, subscriber)
}
