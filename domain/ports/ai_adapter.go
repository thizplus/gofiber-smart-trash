package ports

import (
	"context"
)

// ClassificationResult represents the AI classification response
type ClassificationResult struct {
	Category    string  `json:"category"`
	SubCategory string  `json:"sub_category,omitempty"`
	Confidence  float64 `json:"confidence"`
	BinNumber   int     `json:"bin_number"`
	BinLabel    string  `json:"bin_label"`
	Message     string  `json:"message"`
}

// AIAdapter defines the interface for AI classification service
type AIAdapter interface {
	// ClassifyImage sends an image URL to AI service and returns classification result
	ClassifyImage(ctx context.Context, imageURL string) (*ClassificationResult, error)

	// Health checks if AI service is available
	Health(ctx context.Context) (bool, error)
}
