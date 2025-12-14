package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gofiber-smart-trash/interfaces/api/handlers"
	"gofiber-smart-trash/interfaces/api/routes"
	"gofiber-smart-trash/pkg/di"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize DI container
	container := di.NewContainer()

	// Initialize all dependencies
	if err := container.Initialize(); err != nil {
		log.Fatal("Failed to initialize container:", err)
	}

	// Setup graceful shutdown
	setupGracefulShutdown(container)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: container.GetConfig().App.Name,
	})

	// Create handlers
	h := handlers.NewHandlers(container.GetTrashService())

	// Setup routes (routes include middleware setup)
	routes.SetupRoutes(app, h)

	// Start server
	port := container.GetConfig().App.Port
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌍 Environment: %s", container.GetConfig().App.Env)
	log.Printf("📚 Health check: http://localhost:%s/health", port)
	log.Printf("📖 API endpoints:")
	log.Printf("   GET  /api/upload-url")
	log.Printf("   POST /api/trash")
	log.Printf("   GET  /api/trash")
	log.Printf("   GET  /api/trash/:id")

	log.Fatal(app.Listen(":" + port))
}

func setupGracefulShutdown(container *di.Container) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("\n🛑 Gracefully shutting down...")

		if err := container.Cleanup(); err != nil {
			log.Printf("❌ Error during cleanup: %v", err)
		}

		log.Println("👋 Shutdown complete")
		os.Exit(0)
	}()
}