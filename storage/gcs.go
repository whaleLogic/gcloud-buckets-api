package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// Client wraps Google Cloud Storage client
type Client struct {
	client     *storage.Client
	bucketName string
}

// Ensure Client implements StorageClient interface
var _ StorageClient = (*Client)(nil)

// UploadResult contains the result of an upload operation
type UploadResult struct {
	FileName string `json:"fileName"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
}

// NewClient creates a new storage client
func NewClient(ctx context.Context, projectID, bucketName string) (*Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &Client{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// UploadFile uploads a file to Google Cloud Storage
func (c *Client) UploadFile(ctx context.Context, fileName string, reader io.Reader) (*UploadResult, error) {
	// Create a unique filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	objectName := fmt.Sprintf("%s-%s", timestamp, fileName)

	// Get bucket handle
	bucket := c.client.Bucket(c.bucketName)
	
	// Create object handle
	obj := bucket.Object(objectName)
	
	// Create writer
	writer := obj.NewWriter(ctx)
	defer writer.Close()

	// Copy file content
	size, err := io.Copy(writer, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate public URL
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, objectName)

	return &UploadResult{
		FileName: objectName,
		URL:      url,
		Size:     size,
	}, nil
}

// Close closes the storage client
func (c *Client) Close() error {
	return c.client.Close()
}