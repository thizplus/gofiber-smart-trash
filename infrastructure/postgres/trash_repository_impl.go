package postgres

import (
	"context"

	"gofiber-smart-trash/domain/models"
	"gofiber-smart-trash/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type trashRepositoryImpl struct {
	db *gorm.DB
}

// NewTrashRepository creates a new instance of TrashRepository
func NewTrashRepository(db *gorm.DB) repositories.TrashRepository {
	return &trashRepositoryImpl{db: db}
}

// Create inserts a new trash record into the database
func (r *trashRepositoryImpl) Create(ctx context.Context, trash *models.TrashRecord) error {
	return r.db.WithContext(ctx).Create(trash).Error
}

// FindByID retrieves a trash record by its ID
func (r *trashRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.TrashRecord, error) {
	var trash models.TrashRecord
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&trash).Error; err != nil {
		return nil, err
	}
	return &trash, nil
}

// FindAll retrieves trash records with filtering and pagination
func (r *trashRepositoryImpl) FindAll(ctx context.Context, filter repositories.TrashFilter) ([]models.TrashRecord, int64, error) {
	var trashList []models.TrashRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&models.TrashRecord{})

	// Apply filters
	if filter.DeviceID != "" {
		query = query.Where("device_id = ?", filter.DeviceID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if err := query.
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&trashList).Error; err != nil {
		return nil, 0, err
	}

	return trashList, total, nil
}
