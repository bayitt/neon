package controllers

import (
	"context"
	"mime/multipart"
	"neon/middleware"
	"neon/services"
	"neon/utilities"
	"neon/validators"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReviewController struct {
	validator *validators.ReviewValidator
	service   *services.ReviewService
}

func RegisterReviewRoutes(group *echo.Group) {
	db := utilities.GetDatabaseObject()
	cs := &services.CategoryService{DB: db}
	ss := &services.SeriesService{DB: db}
	rs := &services.ReviewService{DB: db}
	rv := &validators.ReviewValidator{Rs: rs, Cs: cs, Ss: ss}
	rc := &ReviewController{validator: rv, service: rs}

	createReviewGroup := group.Group("")
	createReviewGroup.Use(middleware.AuthMiddleware)
	createReviewGroup.POST("", rc.create)

	updateReviewGroup := group.Group("/:uuid")
	updateReviewGroup.Use(middleware.AuthMiddleware)
	updateReviewGroup.PUT("", rc.update)
}

func uploadImage(image *multipart.FileHeader, uuid string) (string, error) {
	cloudinary, cloudinaryErr := cloudinary.New()
	if cloudinaryErr != nil {
		return "", utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			cloudinaryErr.Error(),
		)
	}

	source, err := image.Open()
	if err != nil {
		return "", utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			err.Error(),
		)
	}

	defer source.Close()
	uploadResult, uploadErr := cloudinary.Upload.Upload(
		context.Background(),
		source,
		uploader.UploadParams{
			Folder:       os.Getenv("CLOUDINARY_FOLDER"),
			ResourceType: "image",
			PublicID:     uuid,
			Overwrite:    api.Bool(true),
		},
	)

	if uploadErr != nil {
		return "", utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			uploadErr.Error(),
		)
	}

	return uploadResult.SecureURL, nil
}

func (rc *ReviewController) create(context echo.Context) error {
	crDto, err := rc.validator.ValidateCreate(context)
	if err != nil {
		return err
	}

	uuid := uuid.New()
	crDto.Uuid = uuid

	imageFile, imageErr := context.FormFile("image")

	if imageErr == nil {
		image, uploadErr := uploadImage(imageFile, uuid.String())

		if uploadErr != nil {
			return uploadErr
		}
		crDto.Image = &image
	}

	review, createErr := rc.service.Create(crDto)
	if createErr != nil {
		return createErr
	}
	return context.JSON(http.StatusCreated, review)
}

func (rc *ReviewController) update(context echo.Context) error {
	review, urDto, err := rc.validator.ValidateUpdate(context)
	if err != nil {
		return err
	}

	imageFile, imageErr := context.FormFile("image")

	if imageErr == nil {
		image, uploadErr := uploadImage(imageFile, review.Uuid.String())

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
