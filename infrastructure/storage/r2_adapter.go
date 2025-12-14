package storage

import (
	"context"
	"fmt"
	"time"

	"gofiber-smart-trash/domain/ports"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2StorageAdapter implements StorageAdapter for Cloudflare R2
type R2StorageAdapter struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

// NewR2StorageAdapter creates a new R2 storage adapter
func NewR2StorageAdapter(accountID, accessKeyID, secretAccessKey, bucket, publicURL string) (ports.StorageAdapter, error) {
	// Cloudflare R2 endpoint format
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	// Load AWS SDK config with R2 credentials
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		)),
		config.WithRegion("auto"), // R2 uses "auto" region
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with R2 endpoint
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &R2StorageAdapter{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}, nil
}

// GeneratePresignedUploadURL generates a presigned URL for uploading objects to R2
func (a *R2StorageAdapter) GeneratePresignedUploadURL(ctx context.Context, key string, expiry time.Duration) (*ports.PresignedURLResponse, error) {
	presignClient := s3.NewPresignClient(a.client)

	// Create a presigned PUT request
	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(key),
		ContentType: aws.String("image/jpeg"),
	}, s3.WithPresignExpires(expiry))

	if err != nil {
		return nil, fmt.Errorf("failed to presign PutObject: %w", err)
	}

	// Generate public URL for accessing the uploaded file
	publicURL := a.GeneratePublicURL(key)

	return &ports.PresignedURLResponse{
		UploadURL: req.URL,
		PublicURL: publicURL,
		ExpiresIn: int64(expiry.Seconds()),
	}, nil
}

// GeneratePublicURL generates a public URL for accessing an object
func (a *R2StorageAdapter) GeneratePublicURL(key string) string {
	return fmt.Sprintf("%s/%s", a.publicURL, key)
}

// DeleteObject removes an object from R2 storage
func (a *R2StorageAdapter) DeleteObject(ctx context.Context, key string) error {
	_, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	})
	return err
}
