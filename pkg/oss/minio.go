package oss

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClientConfig struct {
	Domain          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Prefix          string
}

var _ IClient = (*MinioClient)(nil)

type MinioClient struct {
	config MinioClientConfig
	client *minio.Client
}

func NewMinioClient(config MinioClientConfig) (*MinioClient, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if exists, err := client.BucketExists(ctx, config.BucketName); err != nil {
		return nil, err
	} else if !exists {
		if err := client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &MinioClient{
		config: config,
		client: client,
	}, nil
}

func (c *MinioClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.ReadSeekCloser, objectSize int64, options ...PutObjectOptions) (*PutObjectResult, error) {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	var opt PutobjectOptions
	if len(options) > 0 {
		opt = options[0]
	}

	objectName = formatObjectName(c.config.Prefix, objectName)
	output, err := c.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType:  opt.ContentType,
		UserMetadata: opt.UserMetadata,
	})

	if err != nil {
		return nil, err
	}

	return &PutObjectResult{
		URL:  c.config.Domain + "/" + objectName,
		Key:  output.Key,
		ETag: output.ETag,
		Size: output.Size,
	}, nil

}
