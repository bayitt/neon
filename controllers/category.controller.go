package controllers

import (
	"neon/middleware"
	"neon/services"
	"neon/utilities"
	"neon/validators"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CategoryController struct {
	service   *services.CategoryService
	validator *validators.CategoryValidator
}

func RegisterCategoryRoutes(group *echo.Group) {
	db := utilities.GetDatabaseObject()
	cs := &services.CategoryService{DB: db}
	cc := &CategoryController{service: cs, validator: &validators.CategoryValidator{Service: cs}}

	group.Use(middleware.AuthMiddleware)
	group.POST("", cc.create)
	group.PUT("/:uuid", cc.update)
}

func (cc *CategoryController) create(context echo.Context) error {
	dto, err := cc.validator.ValidateCreate(context)
	if err != nil {
		return err
	}

	category, _ := cc.service.Create(dto)
	return context.JSON(http.StatusCreated, category)
}

func (cc *CategoryController) update(context echo.Context) error {
	category, dto, err := cc.validator.ValidateUpdate(context)
	if err != nil {
		return err
	}

	updatedCategory, _ := cc.service.Update(category, dto)
	return context.JSON(http.StatusOK, updatedCategory)
}
