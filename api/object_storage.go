package api

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func (server *Server) handleFilesUpload(files []*multipart.FileHeader) ([]string, error) {
	if len(files) == 0 {
		return nil, errors.New("no file uploaded")
	}

	// Channels to collect results and errors
	urlsChan := make(chan string, len(files))
	errChan := make(chan error, len(files))

	// Wait group to wait for all Go routines to finish
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()

			url, err := server.uploadFileHandler(file)
			if err != nil {
				errChan <- err
				return
			}
			urlsChan <- url
		}(file)
	}

	// Wait for all uploads to complete
	wg.Wait()
	close(urlsChan)
	close(errChan)

	// Check if there were any errors
	if len(errChan) > 0 {
		return nil, <-errChan // Return the first error
	}

	// Collect all URLs
	urls := make([]string, 0, len(files))
	for url := range urlsChan {
		urls = append(urls, url)
	}

	return urls, nil
}

//func (server *Server) handleFilesUpload(files []*multipart.FileHeader) ([]string, error) {
//	if len(files) == 0 {
//		return nil, errors.New("no file uploaded")
//	}
//
//	urls := make([]string, 0, len(files))
//	for _, file := range files {
//		url, err := server.uploadFileHandler(file)
//		if err != nil {
//			return nil, err
//		}
//		urls = append(urls, url)
//	}
//	return urls, nil
//}

func (server *Server) uploadFileHandler(file *multipart.FileHeader) (string, error) {
	fileName := generateUniqueFilename() + getFileExtension(file)
	multiPartFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(multiPartFile multipart.File) {
		err := multiPartFile.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing file")
		}
	}(multiPartFile)

	return server.uploadFile(multiPartFile, fileName, file.Header.Get("Content-Type"))
}

func (server *Server) uploadFile(file multipart.File, fileName string, contentType string) (string, error) {

	bucketName := server.config.AWSBucketName
	log.Info().Msgf("Uploading file to bucket: %s", bucketName)

	// Upload the file to S3
	_, err := server.svc.PutObject(&s3.PutObjectInput{
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

	return fileName, nil
}

//func (server *Server) uploadFile(file multipart.File, fileName string, contentType string) (string, error) {
//
//	bucketName := server.config.AWSBucketName
//	log.Info().Msgf("Uploading file to bucket: %s", bucketName)
//
//	// Upload the file to S3
//	_, err := server.svc.PutObject(&s3.PutObjectInput{
//		Bucket:               aws.String(bucketName),
//		Key:                  aws.String(fileName),
//		Body:                 file,
//		ContentType:          aws.String(contentType),
//		ContentDisposition:   aws.String("attachment"),
//		ServerSideEncryption: aws.String("AES256"),
//	})
//
//	if err != nil {
//		log.Error().Err(err).Msg("Error uploading file to S3")
//		return "", err
//	}
//
//	url := generateS3URL(server.config.AWSRegion, bucketName, fileName)
//	return url, nil
//}

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
