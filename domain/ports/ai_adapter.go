package ports

import (
	"context"
)

// ClassificationResult represents the AI classification response
type ClassificationResult struct {
	Category     string  `json:"category"`
	SubCategory  string  `json:"sub_category,omitempty"`
	Confidence   float64 `json:"confidence"`
	BinNumber    int     `json:"bin_number"`
	BinLabel     string  `json:"bin_label"`
	Message      string  `json:"message"`
	L0Detected   bool    `json:"l0_detected"`             // L0 พบวัตถุหรือไม่
	L0Label      string  `json:"l0_label,omitempty"`      // YOLO detected object (bottle, cup, etc.)
	L0Confidence float64 `json:"l0_confidence,omitempty"` // YOLO confidence
}

// AIAdapter defines the interface for AI classification service
type AIAdapter interface {
	// ClassifyImage sends an image URL to AI service and returns classification result
	ClassifyImage(ctx context.Context, imageURL string) (*ClassificationResult, error)

	// Health checks if AI service is available
	Health(ctx context.Context) (bool, error)
}
