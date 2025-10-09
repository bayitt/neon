package controllers

import (
	"neon/middleware"
	"neon/services"
	"neon/utilities"
	"neon/validators"
	"net/http"

	"github.com/labstack/echo/v4"
)

type categoryController struct {
	service   *services.CategoryService
	validator *validators.CategoryValidator
}

func RegisterCategoryRoutes(app *echo.Echo) {
	db := utilities.GetDatabaseObject()
	cs := &services.CategoryService{DB: db}
	cc := &categoryController{service: cs, validator: &validators.CategoryValidator{Service: cs}}

	createCategoryGroup := app.Group("/categories")
	createCategoryGroup.Use(middleware.AuthMiddleware)
	createCategoryGroup.POST("", cc.create)

	updateCategoryGroup := app.Group("/categories")
	updateCategoryGroup.Use(middleware.AuthMiddleware)
	updateCategoryGroup.PUT("/:uuid", cc.update)

	app.GET("/categories", cc.get)
}

func (cc *categoryController) create(context echo.Context) error {
	dto, err := cc.validator.ValidateCreate(context)
	if err != nil {
		return err
	}

	category, createErr := cc.service.Create(dto)
	if createErr != nil {
		return err
	}
	return context.JSON(http.StatusCreated, category)
}

func (cc *categoryController) update(context echo.Context) error {
	category, dto, err := cc.validator.ValidateUpdate(context)
	if err != nil {
		return err
	}

	updatedCategory, updateErr := cc.service.Update(category, dto)
	if updateErr != nil {
		return updateErr
	}
	return context.JSON(http.StatusOK, updatedCategory)
}

func (cc *categoryController) get(context echo.Context) error {
	categories, err := cc.service.Find()
	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, categories)
}
