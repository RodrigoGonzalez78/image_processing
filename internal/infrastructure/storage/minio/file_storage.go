package minio

import (
	"bytes"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorage struct {
	client     *minio.Client
	bucketName string
}

func NewFileStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*FileStorage, error) {
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

	return &FileStorage{
		client:     client,
		bucketName: bucket,
	}, nil
}

func (fs *FileStorage) Upload(ctx context.Context, path string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)
	_, err := fs.client.PutObject(ctx, fs.bucketName, path, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (fs *FileStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	return fs.client.GetObject(ctx, fs.bucketName, path, minio.GetObjectOptions{})
}
