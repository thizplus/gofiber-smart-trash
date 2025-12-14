package repositories

import (
	"context"

	"gofiber-smart-trash/domain/models"

	"github.com/google/uuid"
)

type TrashRepository interface {
	Create(ctx context.Context, trash *models.TrashRecord) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.TrashRecord, error)
	FindAll(ctx context.Context, filter TrashFilter) ([]models.TrashRecord, int64, error)
}

type TrashFilter struct {
	DeviceID string
	Limit    int
	Offset   int
}
