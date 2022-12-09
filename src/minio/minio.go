package minio

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	miniov7 "github.com/minio/minio-go/v7"
	zlog "github.com/rs/zerolog/log"
)

type Service struct {
	MinioClient *ObjectStoreClient
}

func (s *Service) GetObject(filePath, fileName, bucketName string) ([]byte, error) {
	objectName := fmt.Sprintf("%s/%s", filePath, fileName)
	objectReader, err := s.MinioClient.GetObject(context.Background(), bucketName, objectName, miniov7.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot download objects from minio %w", err)
	}
	defer objectReader.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(objectReader)
	if err != nil {
		return nil, fmt.Errorf("cannot read object content from reader with: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *Service) PutObject(fileContent []byte, filePath, fileName, bucketName string, retentionPeriod int) error {
	exists, err := s.checkIfBucketExists(bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = s.createNewBucket(bucketName)
		if err != nil {
			return fmt.Errorf("could not create bucket with %w", err)
		}
	}
	conttype := "application/octet-stream"
	size := int64(len(fileContent))
	objectName := fmt.Sprintf("%s/%s", filePath, fileName)
	if len(fileContent) > 512 {
		conttype = http.DetectContentType(fileContent[:512])
	} else {
		conttype = http.DetectContentType(fileContent)
	}
	reader := bytes.NewReader(fileContent)

	_, err = s.MinioClient.PutObject(context.Background(), bucketName, objectName, reader, size, miniov7.PutObjectOptions{ContentType: conttype})
	if err != nil {
		return err
	}
	return nil
}

// CheckIfBucketExists : check if bucket exists
func (s *Service) checkIfBucketExists(bucketName string) (bool, error) {
	// Check to see if we already own this bucket
	exists, errBucketExists := s.MinioClient.BucketExists(context.Background(), bucketName)

	if errBucketExists != nil {
		zlog.Err(errBucketExists)
	}
	return exists, errBucketExists
}

// createNewBucket : Create new bucket
func (s *Service) createNewBucket(bucketName string) error {
	opt := miniov7.MakeBucketOptions{Region: "us-east-1"}
	err := s.MinioClient.MakeBucket(context.Background(), bucketName, opt)

	if err != nil {
		return err
	}
	return nil
}
