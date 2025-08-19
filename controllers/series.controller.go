package controllers

import (
	"neon/middleware"
	"neon/services"
	"neon/utilities"
	"neon/validators"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SeriesController struct {
	service   *services.SeriesService
	validator *validators.SeriesValidator
}

func RegisterSeriesRoutes(group *echo.Group) {
	db := utilities.GetDatabaseObject()
	ss := &services.SeriesService{DB: db}
	sc := &SeriesController{service: ss, validator: &validators.SeriesValidator{Service: ss}}

	group.Use(middleware.AuthMiddleware)
	group.POST("", sc.create)
}

func (sc *SeriesController) create(context echo.Context) error {
	csDto, err := sc.validator.ValidateCreate(context)
	if err != nil {
		return err
	}

	series, _ := sc.service.Create(csDto)
	return context.JSON(http.StatusCreated, series)
}
