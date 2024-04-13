package api

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"strings"
	"time"
)

func (server *Server) uploadfile(file multipart.File, fileName string, contentType string) (string, error) {
	// Create a session with S3
	svc := s3.New(server.awsSession)

	// Upload the file to S3
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(server.config.AWSBucketName),
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

	url := generateS3URL(server.config.AWSRegion, server.config.AWSBucketName, fileName)
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
