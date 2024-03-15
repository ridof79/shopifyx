package delivery

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"shopifyx/util"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	uuid "github.com/nu7hatch/gouuid"
)

const (
	FileRequired                       = "file is required"
	FileFormatNotSupported             = "file format is not supported, only *.jpg or *.jpeg allowed"
	FileSizeExceedsMaximumAllowedSize  = "file size exceeds maximum allowed size of 2MB"
	FileSizeLessThanMinimumAllowedSize = "file size is less than minimum allowed size of 10KB"

	FailedToUploadImage = "failed to upload image"
)

func UploadImageHandler(c echo.Context) error {
	var awsAccesKeyId = os.Getenv("S3_ID")
	var awsSecretAccessKey = os.Getenv("S3_SECRET_KEY")
	var awsBucketName = os.Getenv("S3_BUCKET_NAME")
	var awsRegion = "ap-southeast-1"

	// get file from form data
	file, err := c.FormFile("file")
	if err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, FileRequired)
	}

	// validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" {
		return util.ErrorHandler(c, http.StatusBadRequest, FileFormatNotSupported)
	}

	// validate file size
	if file.Size > 2<<20 {
		return util.ErrorHandler(c, http.StatusBadRequest, FileSizeExceedsMaximumAllowedSize)
	}
	if file.Size < 10<<10 { // 10KB
		return util.ErrorHandler(c, http.StatusBadRequest, FileSizeLessThanMinimumAllowedSize)
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

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)

	// upload file to S3
	uploadResult, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &awsBucketName,
		Key:    &filename,
		Body:   fileContent,
	})

	if err != nil {
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToUploadImage)
	}

	return c.JSON(http.StatusOK,
		map[string]interface{}{
			"image_url": uploadResult.Location,
		})
}
