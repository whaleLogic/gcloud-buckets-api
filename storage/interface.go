package storage

import (
	"context"
	"io"
)

// StorageClient defines the interface for storage operations
type StorageClient interface {
	UploadFile(ctx context.Context, fileName string, reader io.Reader) (*UploadResult, error)
	Close() error
}