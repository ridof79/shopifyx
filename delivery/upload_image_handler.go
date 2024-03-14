package delivery

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	uuid "github.com/nu7hatch/gouuid"
)

func UploadImageHandler(c echo.Context) error {
	var awsAccesKeyId = os.Getenv("S3_ID")
	var awsSecretAccessKey = os.Getenv("S3_SECRET_KEY")
	var awsBucketName = "shopifyx"
	var awsBaseURL = os.Getenv("S3_BASE_URL")
	var awsRegion = "ap-southeast-1"

	// get file from form data
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "failed to get file from form data",
			})
	}

	// validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "file format is not supported, only *.jpg or *.jpeg allowed",
			})
	}

	// validate file size
	if file.Size > 2<<20 {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "file size exceeds maximum allowed size of 2MB",
			})
	}
	if file.Size < 10<<10 { // 10KB
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "file size is less than minimum allowed size of 10KB",
			})
	}
	// generate random filename using UUID
	uuidValue, _ := uuid.NewV4()
	filename := uuidValue.String() + filepath.Ext(file.Filename)

	// Open the file
	fileContent, _ := file.Open()
	defer fileContent.Close()

	// create AWS session

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccesKeyId, awsSecretAccessKey, "")),
		config.WithRegion(awsRegion),
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			})
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &awsBucketName,
		Key:    &filename,
		Body:   fileContent,
	})

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			})
	}

	// create image URL
	imageURL := fmt.Sprintf("https://%s.%s/%s", awsBucketName, awsBaseURL, filename)

	return c.JSON(http.StatusOK,
		map[string]interface{}{
			"image_url": imageURL,
		})
}
