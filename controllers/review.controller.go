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

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type reviewController struct {
	validator *validators.ReviewValidator
	service   *services.ReviewService
}

func RegisterReviewRoutes(app *echo.Echo) {
	db := utilities.GetDatabaseObject()
	cs := &services.CategoryService{DB: db}
	ss := &services.SeriesService{DB: db}
	rs := &services.ReviewService{DB: db}
	rv := &validators.ReviewValidator{Rs: rs, Cs: cs, Ss: ss}
	rc := &reviewController{validator: rv, service: rs}

	createReviewGroup := app.Group("/reviews")
	createReviewGroup.Use(middleware.AuthMiddleware)
	createReviewGroup.POST("", rc.create)

	updateReviewGroup := app.Group("/reviews/:uuid")
	updateReviewGroup.Use(middleware.AuthMiddleware)
	updateReviewGroup.PUT("", rc.update)

	app.GET("/categories/:category_uuid/reviews", rc.getByCategory)
	app.GET("/series/:series_uuid/reviews", rc.getBySeries)
	app.GET("/reviews/:slug", rc.get)
	app.GET("/reviews", rc.getAll)
}

func parseReviews(context echo.Context, reviews []models.Review) []map[string]interface{} {
	query := context.Request().URL.Query()
	fields := query.Get("fields")
	reviewsJson, _ := json.Marshal(reviews)
	var reviewsResponse []map[string]interface{}
	json.Unmarshal(reviewsJson, &reviewsResponse)

	if len(fields) == 0 {
		return reviewsResponse
	}

	parsedFields := strings.Split(fields, ",")
	var parsedReviews = []map[string]interface{}{}

	for _, review := range reviewsResponse {
		var parsedReview = map[string]interface{}{}
		for _, field := range parsedFields {
			parsedReview[field] = review[field]
		}
		parsedReviews = append(parsedReviews, parsedReview)
	}

	return parsedReviews
}

func (rc *reviewController) create(context echo.Context) error {
	crDto, err := rc.validator.ValidateCreate(context)
	if err != nil {
		return err
	}

	uuid := uuid.New()
	crDto.Uuid = uuid

	imageFile, imageErr := context.FormFile("image")

	if imageErr == nil {
		image, uploadErr := utilities.UploadImage(imageFile, uuid.String())

		if uploadErr != nil {
			return utilities.ThrowError(
				http.StatusInternalServerError,
				"REVIEW_003",
				uploadErr.Error(),
			)
		}
		crDto.Image = &image
	}

	review, createErr := rc.service.Create(crDto)
	if createErr != nil {
		return createErr
	}
	return context.JSON(http.StatusCreated, review)
}

func (rc *reviewController) update(context echo.Context) error {
	review, urDto, err := rc.validator.ValidateUpdate(context)
	if err != nil {
		return err
	}

	imageFile, imageErr := context.FormFile("image")

	if imageErr == nil {
		image, uploadErr := utilities.UploadImage(imageFile, review.Uuid.String())

		if uploadErr != nil {
			return uploadErr
		}
		urDto.Image = &image
	}

	updatedReview, updateErr := rc.service.Update(review, urDto)
	if updateErr != nil {
		return updateErr
	}
	return context.JSON(http.StatusOK, updatedReview)
}

func (rc *reviewController) get(context echo.Context) error {
	review, err := rc.validator.ValidateGet(context)
	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, parseReviews(context, []models.Review{review})[0])
}

func (rc *reviewController) getAll(context echo.Context) error {
	grDto, err := rc.validator.ValidateGetMultiple(context)
	if err != nil {
		return err
	}

	offset := (grDto.Page - 1) * grDto.Count
	reviews, getErr := rc.service.Find(offset, grDto.Count, map[string]uint{})
	if getErr != nil {
		return getErr
	}

	totalReviews := rc.service.Count(map[string]uint{})
	totalPages := uint(math.Ceil(float64(totalReviews) / float64(grDto.Count)))
	return context.JSON(
		http.StatusOK,
		map[string]interface{}{
			"reviews": parseReviews(context, reviews),
			"pagination": map[string]uint{
				"currentPage": grDto.Page,
				"totalPages":  totalPages,
			},
		},
	)
}

func (rc *reviewController) getByCategory(context echo.Context) error {
	grbcDto, err := rc.validator.ValidateGetByCategory(context)
	if err != nil {
		return err
	}

	category := grbcDto.Category
	offset := (grbcDto.Page - 1) * grbcDto.Count
	reviews, reviewErr := rc.service.FindCategoryReviews(category, offset, grbcDto.Count)
	if reviewErr != nil {
		return reviewErr
	}

	totalReviews := rc.service.CountCategoryReviews(category)
	totalPages := uint(math.Ceil(float64(totalReviews) / float64(grbcDto.Count)))

	return context.JSON(
		http.StatusOK,
		map[string]interface{}{
			"reviews":    parseReviews(context, reviews),
			"pagination": map[string]uint{"currentPage": grbcDto.Page, "totalPages": totalPages},
		},
	)
}

func (rc *reviewController) getByCategories(context echo.Context) error {
	grbcDto, err := rc.validator.ValidateGetByCategories(context)
	if err != nil {
		return err
	}

}

func (rc *reviewController) getBySeries(context echo.Context) error {
	grbsDto, err := rc.validator.ValidateGetBySeries(context)
	if err != nil {
		return err
	}

	reviews, reviewsErr := rc.service.Find(0, 15, map[string]uint{"series_id": grbsDto.Series.ID})
	if reviewsErr != nil {
		return err
	}

	return context.JSON(http.StatusOK, parseReviews(context, reviews))
}
