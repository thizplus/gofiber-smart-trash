package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gofiber-smart-trash/domain/ports"
)

// ClassifierClient implements AIAdapter interface
type ClassifierClient struct {
	baseURL    string
	httpClient *http.Client
}

// ClassifyRequest is the request body for classification
type ClassifyRequest struct {
	ImageURL string `json:"image_url"`
}

// ClassifyResponse is the response from AI service
type ClassifyResponse struct {
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

// HealthResponse is the response from health endpoint
type HealthResponse struct {
	Status      string `json:"status"`
	ModelLoaded bool   `json:"model_loaded"`
	Device      string `json:"device"`
}

// NewClassifierClient creates a new AI classifier client
func NewClassifierClient(baseURL string, timeout int) *ClassifierClient {
	return &ClassifierClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// ClassifyImage sends an image URL to AI service and returns classification result
func (c *ClassifierClient) ClassifyImage(ctx context.Context, imageURL string) (*ports.ClassificationResult, error) {
	// Prepare request body
	reqBody := ClassifyRequest{ImageURL: imageURL}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/classify", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned status %d", resp.StatusCode)
	}

	// Parse response
	var classifyResp ClassifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&classifyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ports.ClassificationResult{
		Category:     classifyResp.Category,
		SubCategory:  classifyResp.SubCategory,
		Confidence:   classifyResp.Confidence,
		BinNumber:    classifyResp.BinNumber,
		BinLabel:     classifyResp.BinLabel,
		Message:      classifyResp.Message,
		L0Detected:   classifyResp.L0Detected,
		L0Label:      classifyResp.L0Label,
		L0Confidence: classifyResp.L0Confidence,
	}, nil
}

// Health checks if AI service is available
func (c *ClassifierClient) Health(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/health", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("AI service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		return false, fmt.Errorf("failed to decode health response: %w", err)
	}

	return healthResp.Status == "ok" && healthResp.ModelLoaded, nil
}
