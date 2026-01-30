package repository

import (
	"context"
	"io"
)

type FileStorage interface {
	Upload(ctx context.Context, path string, data []byte, contentType string) error

	Get(ctx context.Context, path string) (io.ReadCloser, error)
}
