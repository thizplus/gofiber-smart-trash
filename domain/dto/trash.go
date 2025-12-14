package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs

type CreateTrashRequest struct {
	DeviceID  string  `json:"device_id" validate:"required"`
	ImageURL  string  `json:"image_url" validate:"required,url"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type ListTrashRequest struct {
	DeviceID string `query:"device_id"`
	Limit    int    `query:"limit" validate:"min=1,max=100"`
	Offset   int    `query:"offset" validate:"min=0"`
}

// Response DTOs

type UploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	ImageURL  string `json:"image_url"`
	ExpiresIn int64  `json:"expires_in"`
}

type TrashResponse struct {
	ID        uuid.UUID `json:"id"`
	DeviceID  string    `json:"device_id"`
	ImageURL  string    `json:"image_url"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`

	// Classification results (from AI)
	Category      string    `json:"category"`
	SubCategory   string    `json:"sub_category,omitempty"`
	Confidence    float64   `json:"confidence"`
	BinNumber     int       `json:"bin_number"`
	BinLabel      string    `json:"bin_label"`
	Message       string    `json:"message,omitempty"` // Human-readable result message
	L0Detected    bool      `json:"l0_detected"`              // L0 พบวัตถุหรือไม่
	L0Label       string    `json:"l0_label,omitempty"`       // YOLO detected object (bottle, cup, etc.)
	L0Confidence  float64   `json:"l0_confidence,omitempty"`  // YOLO confidence
	ClassifyError string    `json:"classify_error,omitempty"`
	ClassifiedAt  time.Time `json:"classified_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type ListTrashResponse struct {
	Data       []TrashResponse `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

type Pagination struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
