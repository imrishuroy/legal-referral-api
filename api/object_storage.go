package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

func (server *Server) uploadFile(file multipart.File, fileName string, contentType string, folderName string) (string, error) {
	// Create a session with S3
	svc := s3.New(server.awsSession)

	bucketName := server.config.AWSBucketPrefix + "-" + folderName
	log.Info().Msgf("Uploading file to bucket: %s", bucketName)

	// Upload the file to S3
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fileName),
		Body:                 file,
		ContentType:          aws.String(contentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		log.Error().Err(err).Msg("Error uploading file to S3")
		return "", err
	}

	url := generateS3URL(server.config.AWSRegion, bucketName, fileName)
	return url, nil
}

func preSignS3Object(svc *s3.S3, bucket string, key string) (string, error) {
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	url, err := req.Presign(15 * time.Minute) // Pressing URL for 15 minutes

	if err != nil {
		return "", err
	}
	return url, nil
}

func generateS3URL(region, bucketName, key string) string {
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, key)
	url = strings.ReplaceAll(url, " ", "+")
	return url
}

func getFileExtension(fileHeader *multipart.FileHeader) string {
	// Get the filename from the FileHeader
	filename := fileHeader.Filename

	// Use filepath.Ext to get the extension
	extension := filepath.Ext(filename)

	// Return the extension
	return extension
}

func generateUniqueFilename() string {
	// Get the current time
	timestamp := time.Now().Format("20060102_150405")

	// Generate a random component
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	randomComponent := hex.EncodeToString(randomBytes)

	// Combine timestamp and random component to form the filename
	filename := fmt.Sprintf("%s_%s", timestamp, randomComponent)
	return filename
}

func openFile(fileHeader *multipart.FileHeader) (multipart.File, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	return file, nil
}
