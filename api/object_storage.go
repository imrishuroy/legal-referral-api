package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func (server *Server) handleFilesUpload(files []*multipart.FileHeader) ([]string, error) {
	if len(files) == 0 {
		return nil, errors.New("no file uploaded")
	}

	urls := make([]string, 0, len(files))
	for _, file := range files {
		url, err := uploadFileToS3(file, server.Config.AWSBucketName)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func uploadFileToS3(fileHeader *multipart.FileHeader, bucketName string) (string, error) {
	// Open file from multipart header
	srcFile, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer srcFile.Close()

	// Read the entire file into memory
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, srcFile); err != nil {
		return "", fmt.Errorf("failed to read file into buffer: %w", err)
	}

	key := generateRandomFileName() + getFileExtension(fileHeader)

	contentType := fileHeader.Header.Get("Content-Type")

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to s3: %w", err)
	}

	return key, nil
}

func (srv *Server) uploadFile(file multipart.File, fileName string, contentType string) (string, error) {
	bucketName := srv.Config.AWSBucketName

	_, err := srv.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		log.Error().Err(err).Msg("Error uploading file to S3")
		return "", err
	}

	return fileName, nil
}

//func preSignS3Object(svc *s3.S3, bucket string, key string) (string, error) {
//	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
//		Bucket: aws.String(bucket),
//		Key:    aws.String(key),
//	})
//	url, err := req.Presign(15 * time.Minute) // Pressing URL for 15 minutes
//
//	if err != nil {
//		return "", err
//	}
//	return url, nil
//}

// func generateS3URL(region, bucketName, key string) string {
// 	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, key)
// 	url = strings.ReplaceAll(url, " ", "+")
// 	return url
// }

func getFileExtension(fileHeader *multipart.FileHeader) string {
	// Get the filename from the FileHeader
	filename := fileHeader.Filename

	// Use filepath.Ext to get the extension
	extension := filepath.Ext(filename)

	// Return the extension
	return extension
}

func generateRandomFileName() string {
	bytes := make([]byte, 16) // 16 bytes = 32-character hex string
	_, err := rand.Read(bytes)
	if err != nil {
		log.Error().Err(err).Msg("Error generating random bytes")
		return "default_filename"
	}
	return hex.EncodeToString(bytes)
}

// func openFile(fileHeader *multipart.FileHeader) (multipart.File, error) {
// 	file, err := fileHeader.Open()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return file, nil
// }
