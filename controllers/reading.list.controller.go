package controllers

import (
	"encoding/json"
	"fmt"
	"neon/middleware"
	"neon/models"
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
	guardedGroup.PUT("/:uuid", rlc.update)
}

func parseReadingListItem(readingListItem models.ReadingList) map[string]interface{} {
	readingListItemJson, _ := json.Marshal(readingListItem)
	var parsedReadingListItem map[string]interface{}
	json.Unmarshal(readingListItemJson, &parsedReadingListItem)

	if parsedReadingListItem["status"].(float64) == 1 {
		parsedReadingListItem["status"] = true
	} else {
		parsedReadingListItem["status"] = false
	}

	return parsedReadingListItem
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

	return context.JSON(http.StatusOK, parseReadingListItem(readingListItem))
}

func (rlc *readingListController) update(context echo.Context) error {
	readingListItem, urlDto, err := rlc.validator.ValidateUpdate(context)
	if err != nil {
		return err
	}

	imageFile, imageErr := context.FormFile("image")

	if imageErr == nil {
		image, uploadErr := utilities.UploadImage(imageFile, readingListItem.Uuid.String())

		if uploadErr != nil {
			return utilities.ThrowError(
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				uploadErr.Error(),
			)
		}
		urlDto.Image = &image
	}

	updatedReadingListItem, updateErr := rlc.service.Update(readingListItem, urlDto)
	if updateErr != nil {
		return updateErr
	}

	return context.JSON(http.StatusOK, parseReadingListItem(updatedReadingListItem))
}
