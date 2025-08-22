package validators

import (
	"fmt"
	"neon/dto"
	"neon/models"
	"neon/services"
	"neon/utilities"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type CategoryValidator struct {
	Service *services.CategoryService
}

func (cv *CategoryValidator) ValidateCreate(context echo.Context) (*dto.CreateCategoryDto, error) {
	ccDto := new(dto.CreateCategoryDto)

	if err := context.Bind(ccDto); err != nil {
		return nil, utilities.ThrowError(http.StatusBadRequest, "MALFORMED_REQUEST", err.Error())
	}

	if err := context.Validate(ccDto); err != nil {
		return nil, err
	}

	ccDto.Name = strings.ToLower(ccDto.Name)
	_, err := cv.Service.FindUnique("name", ccDto.Name)

	if err != nil {
		return ccDto, nil
	}
	return nil, utilities.ThrowError(
		http.StatusBadRequest,
		"CATEGORY_001",
		fmt.Sprintf("category with name %s already exists", ccDto.Name),
	)
}

func (cv *CategoryValidator) ValidateUpdate(
	context echo.Context,
) (models.Category, *dto.UpdateCategoryDto, error) {
	ucDto := new(dto.UpdateCategoryDto)

	if err := context.Bind(ucDto); err != nil {
		return models.Category{}, nil, utilities.ThrowError(
			http.StatusBadRequest,
			"MALFORMED_REQUEST",
			err.Error(),
		)
	}

	if err := context.Validate(ucDto); err != nil {
		return models.Category{}, nil, err
	}

	category, err := cv.Service.FindUnique("uuid", ucDto.Uuid.String())
	if err != nil {
		return models.Category{}, nil, utilities.ThrowError(
			http.StatusNotFound,
			"CATEGORY_002",
			fmt.Sprintf("category with uuid %s does not exist", ucDto.Uuid),
		)
	}

	namedCategory, nameErr := cv.Service.FindUnique("name", ucDto.Name)
	if nameErr == nil && namedCategory.Uuid.String() != category.Uuid.String() {
		return models.Category{}, nil, utilities.ThrowError(
			http.StatusBadRequest,
			"CATEGORY_001",
			fmt.Sprintf("category with name %s already exists", ucDto.Name),
		)
	}

	return category, ucDto, nil
}
