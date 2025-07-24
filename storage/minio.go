package storage

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client     *minio.Client
	BucketName string
}

var MinioClientInstance *MinioClient

func StartMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) error {

	var err error
	MinioClientInstance, err = NewMinioClient(endpoint, accessKey, secretKey, bucket, useSSL)

	return err
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinioClient{
		Client:     client,
		BucketName: bucket,
	}, nil
}

func (mc *MinioClient) UploadImage(ctx context.Context, objectName string, file multipart.File, fileSize int64, contentType string) error {
	_, err := mc.Client.PutObject(ctx, mc.BucketName, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (mc *MinioClient) UploadBytes(ctx context.Context, objectName string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)
	_, err := mc.Client.PutObject(ctx, mc.BucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (mc *MinioClient) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	url, err := mc.Client.PresignedGetObject(ctx, mc.BucketName, objectName, expiry, reqParams)
	return url.String(), err
}

func (mc *MinioClient) GetImage(ctx context.Context, objectName string) (io.ReadCloser, error) {
	return mc.Client.GetObject(ctx, mc.BucketName, objectName, minio.GetObjectOptions{})
}
