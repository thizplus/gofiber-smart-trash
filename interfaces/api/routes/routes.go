package routes

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-smart-trash/interfaces/api/handlers"
	"gofiber-smart-trash/interfaces/api/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, h *handlers.Handlers) {
	// Global middleware
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Smart Trash Picker API",
			"version": "1.0.0",
		})
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// API routes
	api := app.Group("/api")

	// Upload URL generation (for presigned URLs)
	api.Get("/upload-url", h.GenerateUploadURL)

	// Trash management routes
	api.Post("/trash", h.CreateTrash)
	api.Get("/trash", h.ListTrash)
	api.Get("/trash/:id", h.GetTrash)
}