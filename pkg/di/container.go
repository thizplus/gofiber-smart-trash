package di

import (
	"log"

	"gofiber-smart-trash/application/services"
	"gofiber-smart-trash/domain/ports"
	domainServices "gofiber-smart-trash/domain/services"
	"gofiber-smart-trash/infrastructure/ai"
	"gofiber-smart-trash/infrastructure/postgres"
	"gofiber-smart-trash/infrastructure/storage"
	"gofiber-smart-trash/pkg/config"

	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	// Configuration
	Config *config.Config

	// Infrastructure
	DB             *gorm.DB
	StorageAdapter ports.StorageAdapter
	AIAdapter      ports.AIAdapter

	// Services
	TrashService domainServices.TrashService
}

func NewContainer() *Container {
	return &Container{}
}

// Initialize sets up all application dependencies
func (c *Container) Initialize() error {
	if err := c.initConfig(); err != nil {
		return err
	}

	if err := c.initDatabase(); err != nil {
		return err
	}

	if err := c.initStorageAdapter(); err != nil {
		return err
	}

	if err := c.initAIAdapter(); err != nil {
		return err
	}

	if err := c.initServices(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	c.Config = cfg
	log.Println("✓ Configuration loaded")
	return nil
}

func (c *Container) initDatabase() error {
	// Initialize Database
	dbConfig := postgres.DatabaseConfig{
		Host:     c.Config.DB.Host,
		Port:     c.Config.DB.Port,
		User:     c.Config.DB.User,
		Password: c.Config.DB.Password,
		DBName:   c.Config.DB.DBName,
		SSLMode:  c.Config.DB.SSLMode,
	}

	db, err := postgres.NewDatabase(dbConfig)
	if err != nil {
		return err
	}
	c.DB = db
	log.Println("✓ Database connected")

	// Run migrations
	if err := postgres.Migrate(db); err != nil {
		return err
	}
	log.Println("✓ Database migrated")

	return nil
}

func (c *Container) initStorageAdapter() error {
	// Initialize storage adapter based on provider
	switch c.Config.Storage.Provider {
	case "r2":
		adapter, err := storage.NewR2StorageAdapter(
			c.Config.Storage.AccountID,
			c.Config.Storage.AccessKeyID,
			c.Config.Storage.SecretAccessKey,
			c.Config.Storage.Bucket,
			c.Config.Storage.PublicURL,
		)
		if err != nil {
			return err
		}
		c.StorageAdapter = adapter
		log.Println("✓ R2 Storage Adapter initialized")

	// Future: Add support for other providers
	// case "s3":
	//     adapter, err := storage.NewS3StorageAdapter(...)
	//     c.StorageAdapter = adapter
	// case "gcs":
	//     adapter, err := storage.NewGCSStorageAdapter(...)
	//     c.StorageAdapter = adapter

	default:
		log.Printf("Warning: Unknown storage provider '%s', using R2 as default", c.Config.Storage.Provider)
		adapter, err := storage.NewR2StorageAdapter(
			c.Config.Storage.AccountID,
			c.Config.Storage.AccessKeyID,
			c.Config.Storage.SecretAccessKey,
			c.Config.Storage.Bucket,
			c.Config.Storage.PublicURL,
		)
		if err != nil {
			return err
		}
		c.StorageAdapter = adapter
		log.Println("✓ R2 Storage Adapter initialized (default)")
	}

	return nil
}

func (c *Container) initAIAdapter() error {
	// Initialize AI adapter for classification service
	c.AIAdapter = ai.NewClassifierClient(
		c.Config.AI.ServiceURL,
		c.Config.AI.Timeout,
	)

	log.Printf("✓ AI Adapter initialized (URL: %s, Timeout: %ds)", c.Config.AI.ServiceURL, c.Config.AI.Timeout)
	return nil
}

func (c *Container) initServices() error {
	// Initialize repository
	trashRepo := postgres.NewTrashRepository(c.DB)

	// Initialize service with repository, storage adapter, and AI adapter
	c.TrashService = services.NewTrashService(trashRepo, c.StorageAdapter, c.AIAdapter)

	log.Println("✓ Services initialized")
	return nil
}

// Cleanup closes all connections and releases resources
func (c *Container) Cleanup() error {
	log.Println("Starting cleanup...")

	// Close database connection
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Warning: Failed to close database connection: %v", err)
			} else {
				log.Println("✓ Database connection closed")
			}
		}
	}

	log.Println("✓ Cleanup completed")
	return nil
}

// GetConfig returns the application configuration
func (c *Container) GetConfig() *config.Config {
	return c.Config
}

// GetTrashService returns the trash service
func (c *Container) GetTrashService() domainServices.TrashService {
	return c.TrashService
}