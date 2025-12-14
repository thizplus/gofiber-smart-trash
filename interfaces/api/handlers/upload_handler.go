package handlers

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-smart-trash/domain/dto"
)

// GenerateUploadURL handles GET /api/upload-url
// Returns a presigned URL for uploading trash images to cloud storage
func (h *Handlers) GenerateUploadURL(c *fiber.Ctx) error {
	// Get device_id from query parameter
	deviceID := c.Query("device_id")
	if deviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Error:   "MISSING_DEVICE_ID",
			Message: "device_id is required",
		})
	}

	// Generate presigned upload URL
	response, err := h.trashService.GenerateUploadURL(c.Context(), deviceID)
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
