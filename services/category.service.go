package services

import (
	"errors"
	"fmt"
	"neon/dto"
	"neon/models"
	"neon/utilities"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type CategoryService struct {
	DB *gorm.DB
}

func (cs *CategoryService) FindUnique(field string, value string) (models.Category, error) {
	var category models.Category
	result := cs.DB.Where(fmt.Sprintf("%s = ?", field), value).First(&category)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return category, utilities.ThrowError(http.StatusNotFound, "CATEGORY_002", fmt.Sprintf("category with field %s and value of %s does not exist", field, value))
	}

	return category, nil
}

func (cs *CategoryService) Create(ccDto *dto.CreateCategoryDto) (models.Category, error) {
	slug := "/" + strings.Replace(ccDto.Name, " ", "-", -1) + "-" + utilities.GenerateRandomString(4)
	category := models.Category{Name: ccDto.Name, Slug: slug}
	result := cs.DB.Save(&category)

	if result.Error != nil {
		return category, utilities.ThrowError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", result.Error.Error())
	}

	return category, nil
}

func (cs *CategoryService) Update(category models.Category, ucDto *dto.UpdateCategoryDto) (models.Category, error) {
	if category.Name == ucDto.Name {
		return category, nil
	}

	slug := "/" + strings.Replace(ucDto.Name, " ", "-", -1) + "-" + utilities.GenerateRandomString(4)
	category.Name = ucDto.Name
	category.Slug = slug
	result := cs.DB.Save(&category)

	if result.Error != nil {
		return category, utilities.ThrowError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", result.Error.Error())
	}

	return category, nil
}
