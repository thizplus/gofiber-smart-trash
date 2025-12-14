package handlers

import (
	"gofiber-smart-trash/domain/services"
)

// Handlers contains all HTTP handlers and services
type Handlers struct {
	trashService services.TrashService
}

// NewHandlers creates a new instance of Handlers with all dependencies
func NewHandlers(trashService services.TrashService) *Handlers {
	return &Handlers{
		trashService: trashService,
	}
}