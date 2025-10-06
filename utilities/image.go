package utilities

import (
	"context"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(image *multipart.FileHeader, uuid string) (string, error) {
	cloudinary, cloudinaryErr := cloudinary.New()
	if cloudinaryErr != nil {
		return "", ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			cloudinaryErr.Error(),
		)
	}

	source, err := image.Open()
	if err != nil {
		return "", ThrowError(
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
		return "", ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			uploadErr.Error(),
		)
	}

	return uploadResult.SecureURL, nil
}
