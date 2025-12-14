package ports

import (
	"context"
	"time"
)

// StorageAdapter defines the interface for cloud storage operations
// This allows us to easily swap between different storage providers (R2, S3, GCS)
// without changing business logic
type StorageAdapter interface {
	// GeneratePresignedUploadURL creates a presigned URL for uploading objects
	GeneratePresignedUploadURL(ctx context.Context, key string, expiry time.Duration) (*PresignedURLResponse, error)

	// GeneratePublicURL returns the public URL for accessing an object
	GeneratePublicURL(key string) string

	// DeleteObject removes an object from storage
	DeleteObject(ctx context.Context, key string) error
}

// PresignedURLResponse contains the URLs for uploading and accessing an object
type PresignedURLResponse struct {
	UploadURL string // Presigned URL for uploading (PUT request)
	PublicURL string // Public URL for accessing the uploaded file
	ExpiresIn int64  // Expiration time in seconds
}
