package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"gofiber-smart-trash/domain/dto"
	"gofiber-smart-trash/domain/models"
	"gofiber-smart-trash/domain/ports"
	"gofiber-smart-trash/domain/repositories"
	"gofiber-smart-trash/domain/services"

	"github.com/google/uuid"
)

type trashServiceImpl struct {
	trashRepo      repositories.TrashRepository
	storageAdapter ports.StorageAdapter
	aiAdapter      ports.AIAdapter
}

// NewTrashService creates a new instance of TrashService
func NewTrashService(trashRepo repositories.TrashRepository, storageAdapter ports.StorageAdapter, aiAdapter ports.AIAdapter) services.TrashService {
	return &trashServiceImpl{
		trashRepo:      trashRepo,
		storageAdapter: storageAdapter,
		aiAdapter:      aiAdapter,
	}
}

// GenerateUploadURL generates a presigned URL for uploading trash images
func (s *trashServiceImpl) GenerateUploadURL(ctx context.Context, deviceID string) (*dto.UploadURLResponse, error) {
	// Generate unique key: trash/{device_id}/{timestamp}.jpg
	timestamp := time.Now().UnixMilli()
	key := fmt.Sprintf("trash/%s/%d.jpg", deviceID, timestamp)

	// Generate presigned URL using storage adapter
	urlResp, err := s.storageAdapter.GeneratePresignedUploadURL(ctx, key, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return &dto.UploadURLResponse{
		UploadURL: urlResp.UploadURL,
		ImageURL:  urlResp.PublicURL,
		ExpiresIn: urlResp.ExpiresIn,
	}, nil
}

// CreateTrashRecord creates a new trash record in the database with AI classification (SYNC mode)
func (s *trashServiceImpl) CreateTrashRecord(ctx context.Context, req *dto.CreateTrashRequest) (*dto.TrashResponse, error) {
	trash := &models.TrashRecord{
		DeviceID:  req.DeviceID,
		ImageURL:  req.ImageURL,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	// SYNC Mode: Call AI service to classify the image before responding
	var classifyResult *ports.ClassificationResult
	var classifyErr error
	var message string

	if s.aiAdapter != nil {
		log.Printf("[AI] Classifying image: %s", req.ImageURL)
		classifyResult, classifyErr = s.aiAdapter.ClassifyImage(ctx, req.ImageURL)

		if classifyErr != nil {
			log.Printf("[AI] Classification failed: %v", classifyErr)
			trash.ClassifyError = classifyErr.Error()
		} else {
			log.Printf("[AI] Classification result: %s (%.2f%%)", classifyResult.Category, classifyResult.Confidence*100)
			trash.Category = classifyResult.Category
			trash.SubCategory = classifyResult.SubCategory
			trash.Confidence = classifyResult.Confidence
			trash.BinNumber = classifyResult.BinNumber
			trash.BinLabel = classifyResult.BinLabel
			trash.ClassifiedAt = time.Now()
			message = classifyResult.Message
		}
	}

	if err := s.trashRepo.Create(ctx, trash); err != nil {
		return nil, fmt.Errorf("failed to create trash record: %w", err)
	}

	return &dto.TrashResponse{
		ID:           trash.ID,
		DeviceID:     trash.DeviceID,
		ImageURL:     trash.ImageURL,
		Latitude:     trash.Latitude,
		Longitude:    trash.Longitude,
		Category:     trash.Category,
		SubCategory:  trash.SubCategory,
		Confidence:   trash.Confidence,
		BinNumber:    trash.BinNumber,
		BinLabel:     trash.BinLabel,
		Message:      message,
		ClassifyError: trash.ClassifyError,
		ClassifiedAt: trash.ClassifiedAt,
		CreatedAt:    trash.CreatedAt,
	}, nil
}

// GetTrashByID retrieves a trash record by its ID
func (s *trashServiceImpl) GetTrashByID(ctx context.Context, id uuid.UUID) (*dto.TrashResponse, error) {
	trash, err := s.trashRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trash record: %w", err)
	}

	return &dto.TrashResponse{
		ID:            trash.ID,
		DeviceID:     trash.DeviceID,
		ImageURL:     trash.ImageURL,
		Latitude:     trash.Latitude,
		Longitude:    trash.Longitude,
		Category:     trash.Category,
		SubCategory:  trash.SubCategory,
		Confidence:   trash.Confidence,
		BinNumber:    trash.BinNumber,
		BinLabel:     trash.BinLabel,
		ClassifyError: trash.ClassifyError,
		ClassifiedAt: trash.ClassifiedAt,
		CreatedAt:    trash.CreatedAt,
	}, nil
}

// ListTrash retrieves a list of trash records with pagination
func (s *trashServiceImpl) ListTrash(ctx context.Context, req *dto.ListTrashRequest) (*dto.ListTrashResponse, error) {
	// Set default values
	if req.Limit == 0 {
		req.Limit = 20
	}

	filter := repositories.TrashFilter{
		DeviceID: req.DeviceID,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}

	trashList, total, err := s.trashRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list trash records: %w", err)
	}

	// Convert to response DTOs
	data := make([]dto.TrashResponse, len(trashList))
	for i, trash := range trashList {
		data[i] = dto.TrashResponse{
			ID:            trash.ID,
			DeviceID:     trash.DeviceID,
			ImageURL:     trash.ImageURL,
			Latitude:     trash.Latitude,
			Longitude:    trash.Longitude,
			Category:     trash.Category,
			SubCategory:  trash.SubCategory,
			Confidence:   trash.Confidence,
			BinNumber:    trash.BinNumber,
			BinLabel:     trash.BinLabel,
			ClassifyError: trash.ClassifyError,
			ClassifiedAt: trash.ClassifiedAt,
			CreatedAt:    trash.CreatedAt,
		}
	}

	return &dto.ListTrashResponse{
		Data: data,
		Pagination: dto.Pagination{
			Total:  total,
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}, nil
}
