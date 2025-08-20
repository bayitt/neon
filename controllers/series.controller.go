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
	group.PUT("/:uuid", sc.update)
}

func (sc *SeriesController) create(context echo.Context) error {
	csDto, err := sc.validator.ValidateCreate(context)
	if err != nil {
		return err
	}

	series, createErr := sc.service.Create(csDto)
	if createErr != nil {
		return createErr
	}
	return context.JSON(http.StatusCreated, series)
}

func (sc *SeriesController) update(context echo.Context) error {
	series, usDto, err := sc.validator.ValidateUpdate(context)
	if err != nil {
		return err
	}

	updatedSeries, updateErr := sc.service.Update(series, usDto)
	if updateErr != nil {
		return updateErr
	}
	return context.JSON(http.StatusOK, updatedSeries)
}
