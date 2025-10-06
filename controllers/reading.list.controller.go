package controllers

import (
	"fmt"
	"neon/middleware"
	"neon/services"
	"neon/utilities"
	"neon/validators"
	"net/http"

	uuidpkg "github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type readingListController struct {
	service   *services.ReadingListService
	validator *validators.ReadingListValidator
}

func RegisterReadingListRoutes(app *echo.Echo) {
	db := utilities.GetDatabaseObject()
	rls := &services.ReadingListService{DB: db}
	rlv := &validators.ReadingListValidator{Service: rls}
	rlc := &readingListController{service: rls, validator: rlv}

	guardedGroup := app.Group("/reading-list")
	guardedGroup.Use(middleware.AuthMiddleware)
	guardedGroup.POST("", rlc.create)
}

func (rlc *readingListController) create(context echo.Context) error {
	crlDto, err := rlc.validator.ValidateCreate(context)
	fmt.Println(err)
	if err != nil {
		return err
	}

	uuid := uuidpkg.New()
	crlDto.Uuid = uuid

	imageFile, imageErr := context.FormFile("image")

	if imageErr != nil {
		return utilities.ThrowError(
			http.StatusBadRequest,
			"READING_LIST_003",
			fmt.Sprintf("Image is missing - %s", imageErr.Error()),
		)
	}

	image, uploadErr := utilities.UploadImage(imageFile, uuid.String())

	if uploadErr != nil {
		return utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			uploadErr.Error(),
		)
	}
	crlDto.Image = image

	readingListItem, createErr := rlc.service.Create(crlDto)

	if createErr != nil {
		return createErr
	}

	return context.JSON(http.StatusCreated, readingListItem)
}

func (rlc *readingListController) update(context echo.Context) error {

}
