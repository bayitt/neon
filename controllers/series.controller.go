package controllers

import (
	"encoding/json"
	"math"
	"neon/middleware"
	"neon/models"
	"neon/services"
	"neon/utilities"
	"neon/validators"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type seriesController struct {
	service   *services.SeriesService
	validator *validators.SeriesValidator
}

func RegisterSeriesRoutes(app *echo.Echo) {
	db := utilities.GetDatabaseObject()
	ss := &services.SeriesService{DB: db}
	sc := &seriesController{service: ss, validator: &validators.SeriesValidator{Service: ss}}

	guardedGroup := app.Group("/series")
	guardedGroup.Use(middleware.AuthMiddleware)
	guardedGroup.POST("", sc.create)
	guardedGroup.PUT("/:uuid", sc.update)

	unGuardedGroup := app.Group("/series")
	unGuardedGroup.GET("", sc.getAll)
	unGuardedGroup.GET("/:slug", sc.getBySlug)
}

func parseSeries(context echo.Context, series []models.Series) []map[string]interface{} {
	query := context.Request().URL.Query()
	fields := query.Get("fields")
	seriesJson, _ := json.Marshal(series)
	var seriesResponse []map[string]interface{}
	json.Unmarshal(seriesJson, &seriesResponse)

	if len(fields) == 0 {
		return seriesResponse
	}

	parsedFields := strings.Split(fields, ",")
	var parsedSeries = []map[string]interface{}{}

	for _, series := range seriesResponse {
		var parsedSerie = map[string]interface{}{}

		for _, field := range parsedFields {
			parsedSerie[field] = series[field]
		}

		parsedSeries = append(parsedSeries, parsedSerie)
	}

	return parsedSeries
}

func (sc *seriesController) create(context echo.Context) error {
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

func (sc *seriesController) update(context echo.Context) error {
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

func (sc *seriesController) getBySlug(context echo.Context) error {
	series, err := sc.validator.ValidateGetBySlug(context)
	if err != nil {
		return err
	}

	parsedSeries := parseSeries(context, []models.Series{series})[0]
	return context.JSON(http.StatusOK, parsedSeries)
}

func (sc *seriesController) getAll(context echo.Context) error {
	gsDto, err := sc.validator.ValidateFind(context)
	if err != nil {
		return err
	}

	offset := (gsDto.Page - 1) * gsDto.Count
	series, seriesErr := sc.service.Find(offset, gsDto.Count)
	if seriesErr != nil {
		return seriesErr
	}

	totalSeries := sc.service.Count()
	totalPages := uint(math.Ceil(float64(totalSeries) / float64(gsDto.Count)))

	return context.JSON(
		http.StatusOK,
		map[string]interface{}{
			"series":     parseSeries(context, series),
			"pagination": map[string]uint{"currentPage": gsDto.Page, "totalPages": totalPages},
		},
	)
}
