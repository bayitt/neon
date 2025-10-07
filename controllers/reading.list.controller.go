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
	guardedGroup.DELETE("/:uuid", rlc.delete)

	unguardedGroup := app.Group("/reading-list")
	unguardedGroup.GET("", rlc.getAll)
}

func parseReadingListItems(readingListItem []models.ReadingList) []map[string]interface{} {
	readingListItemJson, _ := json.Marshal(readingListItem)
	var readingListItems []map[string]interface{}
	json.Unmarshal(readingListItemJson, &readingListItems)
	var parsedReadingListItems = []map[string]interface{}{}

	for _, readingListItem := range readingListItems {
		if readingListItem["status"].(float64) == 1 {
			readingListItem["status"] = true
		} else {
			readingListItem["status"] = false
		}

		parsedReadingListItems = append(parsedReadingListItems, readingListItem)
	}

	return parsedReadingListItems
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

	return context.JSON(
		http.StatusOK,
		parseReadingListItems([]models.ReadingList{readingListItem})[0],
	)
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

	return context.JSON(
		http.StatusOK,
		parseReadingListItems([]models.ReadingList{updatedReadingListItem})[0],
	)
}

func (rlc *readingListController) delete(context echo.Context) error {
	readingListItem, err := rlc.validator.ValidateDelete(context)
	if err != nil {
		return err
	}

	deleteErr := rlc.service.Delete(readingListItem)
	if deleteErr != nil {
		return deleteErr
	}

	return context.JSON(http.StatusNoContent, map[string]string{})
}

func (rlc *readingListController) getAll(context echo.Context) error {
	readingListItems, err := rlc.service.Find()

	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, parseReadingListItems(readingListItems))
}
