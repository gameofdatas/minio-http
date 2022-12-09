package minio

import (
	"context"
	"io"

	"github.com/minio-rest/config"
	miniov7 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStoreClient struct {
	Client ObjectStore
}

//go:generate mockgen -source=src/minio/client.go -package minio -destination src/minio/minio_mock.go
type ObjectStore interface {
	MakeBucket(ctx context.Context, bucketName string, opts miniov7.MakeBucketOptions) error
	GetObject(ctx context.Context, bucketName, objectName string, opts miniov7.GetObjectOptions) (*miniov7.Object, error)
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
		opts miniov7.PutObjectOptions) (info miniov7.UploadInfo, err error)
	BucketExists(ctx context.Context, bucketName string) (bool, error)
}

func (o *ObjectStoreClient) MakeBucket(ctx context.Context, bucketName string, opts miniov7.MakeBucketOptions) error {
	return o.Client.MakeBucket(ctx, bucketName, opts)
}

func (o *ObjectStoreClient) GetObject(ctx context.Context, bucketName, objectName string, opts miniov7.GetObjectOptions) (*miniov7.Object, error) {
	return o.Client.GetObject(ctx, bucketName, objectName, opts)
}

func (o *ObjectStoreClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts miniov7.PutObjectOptions) (info miniov7.UploadInfo, err error) {
	return o.Client.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}

func (o *ObjectStoreClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return o.Client.BucketExists(ctx, bucketName)
}

func NewClient(conf config.Provider) (ObjectStore, error) {
	minioClient, err := miniov7.New(conf.GetString("MINIO_ENDPOINT"), &miniov7.Options{
		Creds:  credentials.NewStaticV4(conf.GetString("MINIO_ACCESS_KEY"), conf.GetString("MINIO_SECRET_KEY"), ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &ObjectStoreClient{
		Client: minioClient,
	}, nil
}
