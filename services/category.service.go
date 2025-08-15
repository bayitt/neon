package services

import (
	"errors"
	"fmt"
	"math/rand"
	"neon/dto"
	"neon/models"
	"strings"

	"gorm.io/gorm"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}

type CategoryService struct {
	DB *gorm.DB
}

func (cs *CategoryService) FindUnique(field string, value string) (models.Category, error) {
	var category models.Category
	result := cs.DB.Where(fmt.Sprintf("%s = ?", field), value).First(&category)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return category, fmt.Errorf("category with field %s and value of %s does not exist", field, value)
	}

	return category, nil
}

func (cs *CategoryService) Create(ccDto *dto.CreateCategoryDto) (*models.Category, error) {
	slug := "/" + strings.Replace(ccDto.Name, " ", "-", -1) + "-" + generateRandomString(4)
	category := &models.Category{Name: ccDto.Name, Slug: slug}
	result := cs.DB.Save(&category)

	if result.Error != nil {
		return category, fmt.Errorf("there was an issue creating the category")
	}

	return category, nil
}

func (cs *CategoryService) Update(category models.Category, ucDto *dto.UpdateCategoryDto) (models.Category, error) {
	if category.Name == ucDto.Name {
		return category, nil
	}

	slug := "/" + strings.Replace(ucDto.Name, " ", "-", -1) + "-" + generateRandomString(4)
	category.Name = ucDto.Name
	category.Slug = slug
	result := cs.DB.Save(&category)

	if result.Error != nil {
		return category, fmt.Errorf(result.Error.Error())
	}

	return category, nil
}
