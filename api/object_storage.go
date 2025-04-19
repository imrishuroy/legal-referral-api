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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

func (srv *Server) uploadFileHandler(file *multipart.FileHeader) (string, error) {
	// fileName := generateRandomFileName() + getFileExtension(file)
	multiPartFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer multiPartFile.Close()

	return srv.uploadFile(multiPartFile, file.Filename, file.Header.Get("Content-Type"))
}

func (srv *Server) uploadFile(file multipart.File, fileName string, contentType string) (string, error) {
	bucketName := srv.Config.AWSBucketName
	log.Info().Msgf("Uploading file to bucket: %s", bucketName)
	log.Info().Msgf("File name: %s", fileName)
	log.Info().Msgf("Content type: %s", contentType)

	_, err := srv.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("legalreferral-media"),
		Key:    aws.String(fileName),
		Body:   file,
		// ContentLength: &contentLength,
		ContentType: aws.String(contentType),
		// ACL: types.ObjectCannedACLPublicRead,
		// ContentDisposition:   aws.String("attachment"),
		// ServerSideEncryption: types.ServerSideEncryptionAes256,
	})

	if err != nil {
		log.Error().Err(err).Msg("Error uploading file to S3")
		return "", err
	}

	return fileName, nil
}

func (srv *Server) uploadFile2(file multipart.File, fileName string, contentType string) (string, error) {
	bucketName := srv.Config.AWSBucketName
	log.Info().Msgf("Uploading file to bucket: %s", bucketName)
	log.Info().Msgf("File name: %s", fileName)
	log.Info().Msgf("Content type: %s", contentType)

	// Seek back to beginning if possible
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			log.Error().Err(err).Msg("Failed to seek to start of file before upload")
			return "", err
		}
	} else {
		log.Warn().Msg("File is not seekable")
	}

	_, err := srv.S3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,

		ContentType:          aws.String(contentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	})

	if err != nil {
		log.Error().Err(err).Msg("Error uploading file to S3")
		return "", err
	}

	return fileName, nil
}

func (srv *Server) uploadFile3(ctx context.Context, file multipart.File, fileName string, contentType string) (string, error) {

	bucketName := srv.Config.AWSBucketName
	log.Info().Msgf("Uploading file to bucket: %s", bucketName)
	log.Info().Msgf("Uploading file to bucket: %s", bucketName)
	log.Info().Msgf("File name: %s", fileName)
	log.Info().Msgf("Content type: %s", contentType)

	_, err := srv.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(fileName),
		Body:                 file,
		ContentType:          aws.String(contentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: types.ServerSideEncryption(*aws.String("AES256")),
	})

	if err != nil {
		// handle EntityTooLarge error
		// if s3Err, ok := err.(*types.InvalidObjectState); ok {
		// 	if s3Err.ErrorCode() == "EntityTooLarge" {
		// 		log.Error().Msg("File size exceeds the limit")
		// 		return "", errors.New("file size exceeds the limit")
		// 	}
		// }

		log.Error().Err(err).Msg("Error uploading file to S3")
		return "", err
	}

	return fileName, nil
}

// Upload the file to S3
//_, err := s.s.PutObject(&s3.PutObjectInput{
//	Bucket:               aws.String(bucketName),
//	Key:                  aws.String(fileName),
//	Body:                 file,
//	ContentType:          aws.String(contentType),
//	ContentDisposition:   aws.String("attachment"),
//	ServerSideEncryption: aws.String("AES256"),
//})

//func (server *Server) uploadFile(file multipart.File, fileName string, contentType string) (string, error) {
//
//	bucketName := server.Config.AWSBucketName
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
//	url := generateS3URL(server.Config.AWSRegion, bucketName, fileName)
//	return url, nil
//}

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
