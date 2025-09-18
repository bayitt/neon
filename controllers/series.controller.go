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

func RegisterSeriesRoutes(app *echo.Echo) {
	db := utilities.GetDatabaseObject()
	ss := &services.SeriesService{DB: db}
	sc := &SeriesController{service: ss, validator: &validators.SeriesValidator{Service: ss}}

	guardedGroup := app.Group("/series")
	guardedGroup.Use(middleware.AuthMiddleware)
	guardedGroup.POST("", sc.create)
	guardedGroup.PUT("/:uuid", sc.update)

	unGuardedGroup := app.Group("/series")
	unGuardedGroup.GET("", sc.getAll)
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

func (sc *SeriesController) getAll(context echo.Context) error {
	gsDto, err := sc.validator.ValidateFind(context)
	if err != nil {
		return err
	}

	offset := (gsDto.Page - 1) * gsDto.Count
	series, seriesErr := sc.service.Find(offset, gsDto.Count)
	if seriesErr != nil {
		return seriesErr
	}

	return context.JSON(http.StatusOK, series)
}
