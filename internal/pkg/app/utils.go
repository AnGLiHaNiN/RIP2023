package app

import (
	"fmt"
	"path/filepath"


	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func (app *Application) uploadImage(c *gin.Context, image *multipart.FileHeader, UUID string) (*string, error) {
	src, err := image.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	extension := filepath.Ext(image.Filename)
	if extension != ".jpg" && extension != ".jpeg" {
		return nil, fmt.Errorf("разрешены только jpeg изображения")
	}
	imageName := UUID + extension

	_, err = app.minioClient.PutObject(c, app.config.BucketName, imageName, src, image.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return nil, err
	}
	imageURL := fmt.Sprintf("%s/%s/%s", app.config.MinioEndpoint, app.config.BucketName, imageName)
	return &imageURL, nil
}

func (app *Application) deleteImage(c *gin.Context, UUID string) error {
	imageName := UUID + ".jpg"
	fmt.Println(imageName)
	err := app.minioClient.RemoveObject(c, app.config.BucketName, imageName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) getCustomer() string {
	return "63f6cb81-890d-42ff-ba83-721740b5f57e"
}

func (app *Application) getModerator() *string {
	moderaorId := "079b1514-c4c0-4a1d-966d-6c9135a962b2"
	return &moderaorId
}