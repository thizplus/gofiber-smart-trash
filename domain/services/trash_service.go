package services

import (
	"context"

	"gofiber-smart-trash/domain/dto"

	"github.com/google/uuid"
)

type TrashService interface {
	GenerateUploadURL(ctx context.Context, deviceID string) (*dto.UploadURLResponse, error)
	CreateTrashRecord(ctx context.Context, req *dto.CreateTrashRequest) (*dto.TrashResponse, error)
	GetTrashByID(ctx context.Context, id uuid.UUID) (*dto.TrashResponse, error)
	ListTrash(ctx context.Context, req *dto.ListTrashRequest) (*dto.ListTrashResponse, error)
}
