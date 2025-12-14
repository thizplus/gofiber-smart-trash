package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gofiber-smart-trash/domain/dto"
	"gofiber-smart-trash/pkg/utils"
)

// CreateTrash handles POST /api/trash
// Creates a new trash record after image has been uploaded
func (h *Handlers) CreateTrash(c *fiber.Ctx) error {
	var req dto.CreateTrashRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Error:   "VALIDATION_ERROR",
			Message: err.Error(),
		})
	}

	// Create trash record
	response, err := h.trashService.CreateTrashRecord(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Error:   "INTERNAL_ERROR",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetTrash handles GET /api/trash/:id
// Retrieves a single trash record by ID
func (h *Handlers) GetTrash(c *fiber.Ctx) error {
	// Parse UUID from URL parameter
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Error:   "INVALID_ID",
			Message: "Invalid UUID format",
		})
	}

	// Get trash record
	response, err := h.trashService.GetTrashByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Error:   "NOT_FOUND",
			Message: "Trash record not found",
		})
	}

	return c.JSON(dto.APIResponse{
		Success: true,
		Data:    response,
	})
}

// ListTrash handles GET /api/trash
// Retrieves a list of trash records with optional filtering and pagination
func (h *Handlers) ListTrash(c *fiber.Ctx) error {
	var req dto.ListTrashRequest

	// Parse query parameters
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
	}

	// List trash records
	response, err := h.trashService.ListTrash(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Error:   "INTERNAL_ERROR",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.APIResponse{
		Success: true,
		Data:    response,
	})
}
