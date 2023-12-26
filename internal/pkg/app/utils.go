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
	return "ec62488b-29eb-4484-be7a-1ef79454a2a0"
}

func (app *Application) getModerator() *string {
	moderaorId := "dabad2aa-e7e2-4947-87a2-9fbbdb05a860"
	return &moderaorId
}
