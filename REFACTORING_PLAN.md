# 🔄 Smart Trash Picker - Refactoring Plan

## 📋 Overview

ปรับโครงสร้างจาก Clean Architecture แบบเต็มรูปแบบ (User/Auth/Task/Job/File) ให้เหลือเฉพาะส่วนที่จำเป็นสำหรับ **Smart Trash Picker** โดยใช้ **Adapter Pattern** สำหรับ Storage

---

## 🎯 Goals

1. ✅ ลบส่วนที่ไม่เกี่ยวข้อง (Auth, User, Task, Job, WebSocket)
2. ✅ เหลือเฉพาะ Trash Record management
3. ✅ ใช้ Adapter Pattern สำหรับ Storage (เพื่อให้เปลี่ยน provider ได้ง่าย)
4. ✅ Clean Architecture ที่เรียบง่ายและตรงจุด

---

## 🏗️ New Architecture

### Clean Architecture Layers

```
┌────────────────────────────────────────────────────────────┐
│                    Interfaces Layer                        │
│                  (API Handlers & Routes)                   │
│                                                            │
│  GET  /api/upload-url                                      │
│  POST /api/trash                                           │
│  GET  /api/trash                                           │
│  GET  /api/trash/:id                                       │
└────────────────────┬───────────────────────────────────────┘
                     │
┌────────────────────▼───────────────────────────────────────┐
│                   Application Layer                        │
│                   (Business Logic)                         │
│                                                            │
│  TrashService                                              │
│  - CreateTrash(dto) → TrashRecord                          │
│  - GetTrashByID(id) → TrashRecord                          │
│  - ListTrash(filter) → []TrashRecord                       │
│  - GenerateUploadURL(deviceID) → UploadURLResponse         │
└────────────────────┬───────────────────────────────────────┘
                     │
┌────────────────────▼───────────────────────────────────────┐
│                     Domain Layer                           │
│                 (Core Business Objects)                    │
│                                                            │
│  Models:                                                   │
│  - TrashRecord                                             │
│                                                            │
│  Repositories (Interfaces):                                │
│  - TrashRepository                                         │
│                                                            │
│  Storage Ports (Interfaces):                               │
│  - StorageAdapter                                          │
└────────────────────┬───────────────────────────────────────┘
                     │
┌────────────────────▼───────────────────────────────────────┐
│                 Infrastructure Layer                       │
│              (External Services & DB)                      │
│                                                            │
│  PostgreSQL:                      Storage Adapters:        │
│  - TrashRepositoryImpl            - R2StorageAdapter       │
│                                   - S3StorageAdapter        │
│                                   - GCSStorageAdapter       │
└────────────────────────────────────────────────────────────┘
```

---

## 🎨 Adapter Pattern for Storage

### Storage Interface (Port)

```go
// domain/ports/storage_adapter.go
type StorageAdapter interface {
    GeneratePresignedUploadURL(ctx context.Context, key string, expiry time.Duration) (PresignedURLResponse, error)
    GeneratePublicURL(key string) string
    DeleteObject(ctx context.Context, key string) error
}

type PresignedURLResponse struct {
    UploadURL string
    PublicURL string
    ExpiresIn int64
}
```

### Adapter Implementations

```go
// infrastructure/storage/r2_adapter.go
type R2StorageAdapter struct {
    client    *s3.Client
    bucket    string
    publicURL string
}

// infrastructure/storage/s3_adapter.go (Future)
type S3StorageAdapter struct {
    client    *s3.Client
    bucket    string
    region    string
}

// infrastructure/storage/gcs_adapter.go (Future)
type GCSStorageAdapter struct {
    client    *storage.Client
    bucket    string
}
```

### Configuration

```go
// pkg/config/config.go
type StorageConfig struct {
    Provider  string // "r2", "s3", "gcs"
    Bucket    string
    PublicURL string
    // R2/S3 specific
    AccountID       string
    AccessKeyID     string
    SecretAccessKey string
    Region          string
    // GCS specific
    ProjectID       string
    CredentialsPath string
}
```

---

## 📁 New Project Structure

```
gofiber-smart-trash/
│
├── cmd/
│   └── api/
│       └── main.go                      # Entry point
│
├── domain/                              # Core Domain
│   ├── models/
│   │   └── trash.go                     # TrashRecord model
│   │
│   ├── repositories/
│   │   └── trash_repository.go          # TrashRepository interface
│   │
│   ├── services/
│   │   └── trash_service.go             # TrashService interface
│   │
│   ├── ports/
│   │   └── storage_adapter.go           # StorageAdapter interface (Port)
│   │
│   └── dto/
│       ├── trash.go                     # DTOs for trash operations
│       └── common.go                    # Common response structures
│
├── application/                         # Application Layer
│   └── services/
│       └── trash_service_impl.go        # TrashService implementation
│
├── infrastructure/                      # Infrastructure Layer
│   ├── postgres/
│   │   ├── database.go                  # DB connection & migration
│   │   └── trash_repository_impl.go     # TrashRepository implementation
│   │
│   └── storage/                         # Storage Adapters
│       ├── r2_adapter.go                # Cloudflare R2 implementation
│       ├── s3_adapter.go                # AWS S3 (future)
│       └── gcs_adapter.go               # Google Cloud Storage (future)
│
├── interfaces/                          # Interface Layer
│   └── api/
│       ├── handlers/
│       │   ├── handlers.go              # Handler container
│       │   ├── upload_handler.go        # GET /api/upload-url
│       │   └── trash_handler.go         # Trash CRUD handlers
│       │
│       ├── middleware/
│       │   ├── cors_middleware.go       # CORS
│       │   └── logger_middleware.go     # Logging
│       │
│       └── routes/
│           └── routes.go                # Route definitions
│
├── pkg/                                 # Shared Utilities
│   ├── config/
│   │   └── config.go                    # Configuration management
│   │
│   ├── di/
│   │   └── container.go                 # Dependency Injection
│   │
│   └── utils/
│       ├── response.go                  # HTTP response helpers
│       └── validator.go                 # Validation utilities
│
├── .env                                 # Environment variables
├── .env.example                         # Environment template
├── .air.toml                            # Air hot reload config
├── docker-compose.yml                   # Docker orchestration
├── Dockerfile                           # Container definition
├── Makefile                             # Build commands
├── go.mod
├── go.sum
│
├── smart-trash-picker-plan-v3.md        # Original plan
└── REFACTORING_PLAN.md                  # This file
```

---

## 🗑️ Files to DELETE

### Complete Directories
```bash
rm -rf domain/models/user.go
rm -rf domain/models/task.go
rm -rf domain/models/file.go
rm -rf domain/models/job.go

rm -rf domain/repositories/user_repository.go
rm -rf domain/repositories/task_repository.go
rm -rf domain/repositories/file_repository.go
rm -rf domain/repositories/job_repository.go

rm -rf domain/services/user_service.go
rm -rf domain/services/task_service.go
rm -rf domain/services/file_service.go
rm -rf domain/services/job_service.go

rm -rf domain/dto/auth.go
rm -rf domain/dto/user.go
rm -rf domain/dto/task.go
rm -rf domain/dto/file.go
rm -rf domain/dto/job.go

rm -rf application/serviceimpl/user_service_impl.go
rm -rf application/serviceimpl/task_service_impl.go
rm -rf application/serviceimpl/file_service_impl.go
rm -rf application/serviceimpl/job_service_impl.go

rm -rf infrastructure/postgres/user_repository_impl.go
rm -rf infrastructure/postgres/task_repository_impl.go
rm -rf infrastructure/postgres/file_repository_impl.go
rm -rf infrastructure/postgres/job_repository_impl.go

rm -rf infrastructure/redis/
rm -rf infrastructure/websocket/
rm -rf infrastructure/storage/bunny_storage.go

rm -rf interfaces/api/handlers/user_handler.go
rm -rf interfaces/api/handlers/task_handler.go
rm -rf interfaces/api/handlers/file_handler.go
rm -rf interfaces/api/handlers/job_handler.go

rm -rf interfaces/api/middleware/auth_middleware.go
rm -rf interfaces/api/middleware/error_middleware.go

rm -rf interfaces/api/routes/auth_routes.go
rm -rf interfaces/api/routes/user_routes.go
rm -rf interfaces/api/routes/task_routes.go
rm -rf interfaces/api/routes/file_routes.go
rm -rf interfaces/api/routes/job_routes.go
rm -rf interfaces/api/routes/health_routes.go
rm -rf interfaces/api/routes/websocket_routes.go

rm -rf interfaces/api/websocket/

rm -rf pkg/scheduler/
rm -rf pkg/utils/jwt.go
rm -rf pkg/utils/path.go
```

### Directories to Remove
- `application/serviceimpl/` → เปลี่ยนเป็น `application/services/`
- `infrastructure/redis/` → ไม่ใช้
- `infrastructure/websocket/` → ไม่ใช้
- `infrastructure/storage/bunny_storage.go` → เปลี่ยนเป็น adapter pattern
- `interfaces/api/websocket/` → ไม่ใช้
- `pkg/scheduler/` → ไม่ใช้

---

## ✨ Files to CREATE

### 1. Domain Layer

#### `domain/models/trash.go`
```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type TrashRecord struct {
    ID         uint            `gorm:"primaryKey" json:"id"`
    DeviceID   string          `gorm:"type:varchar(20);not null;index" json:"device_id"`
    ImageURL   string          `gorm:"type:text;not null" json:"image_url"`
    Latitude   float64         `gorm:"type:decimal(10,8);not null" json:"latitude"`
    Longitude  float64         `gorm:"type:decimal(11,8);not null" json:"longitude"`
    CreatedAt  time.Time       `json:"created_at"`
    UpdatedAt  time.Time       `json:"updated_at"`
    DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (TrashRecord) TableName() string {
    return "trash_records"
}
```

#### `domain/repositories/trash_repository.go`
```go
package repositories

import (
    "context"
    "gofiber-smart-trash/domain/models"
)

type TrashRepository interface {
    Create(ctx context.Context, trash *models.TrashRecord) error
    FindByID(ctx context.Context, id uint) (*models.TrashRecord, error)
    FindAll(ctx context.Context, filter TrashFilter) ([]models.TrashRecord, int64, error)
}

type TrashFilter struct {
    DeviceID string
    Limit    int
    Offset   int
}
```

#### `domain/services/trash_service.go`
```go
package services

import (
    "context"
    "gofiber-smart-trash/domain/dto"
)

type TrashService interface {
    GenerateUploadURL(ctx context.Context, deviceID string) (*dto.UploadURLResponse, error)
    CreateTrashRecord(ctx context.Context, req *dto.CreateTrashRequest) (*dto.TrashResponse, error)
    GetTrashByID(ctx context.Context, id uint) (*dto.TrashResponse, error)
    ListTrash(ctx context.Context, req *dto.ListTrashRequest) (*dto.ListTrashResponse, error)
}
```

#### `domain/ports/storage_adapter.go`
```go
package ports

import (
    "context"
    "time"
)

type StorageAdapter interface {
    GeneratePresignedUploadURL(ctx context.Context, key string, expiry time.Duration) (*PresignedURLResponse, error)
    GeneratePublicURL(key string) string
    DeleteObject(ctx context.Context, key string) error
}

type PresignedURLResponse struct {
    UploadURL string
    PublicURL string
    ExpiresIn int64
}
```

#### `domain/dto/trash.go`
```go
package dto

import "time"

// Request DTOs
type CreateTrashRequest struct {
    DeviceID  string  `json:"device_id" validate:"required"`
    ImageURL  string  `json:"image_url" validate:"required,url"`
    Latitude  float64 `json:"latitude" validate:"required,latitude"`
    Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type ListTrashRequest struct {
    DeviceID string `query:"device_id"`
    Limit    int    `query:"limit" validate:"min=1,max=100"`
    Offset   int    `query:"offset" validate:"min=0"`
}

// Response DTOs
type UploadURLResponse struct {
    UploadURL string `json:"upload_url"`
    ImageURL  string `json:"image_url"`
    ExpiresIn int64  `json:"expires_in"`
}

type TrashResponse struct {
    ID        uint      `json:"id"`
    DeviceID  string    `json:"device_id"`
    ImageURL  string    `json:"image_url"`
    Latitude  float64   `json:"latitude"`
    Longitude float64   `json:"longitude"`
    CreatedAt time.Time `json:"created_at"`
}

type ListTrashResponse struct {
    Data       []TrashResponse `json:"data"`
    Pagination Pagination      `json:"pagination"`
}

type Pagination struct {
    Total  int64 `json:"total"`
    Limit  int   `json:"limit"`
    Offset int   `json:"offset"`
}
```

#### `domain/dto/common.go`
```go
package dto

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Message string      `json:"message,omitempty"`
}
```

### 2. Application Layer

#### `application/services/trash_service_impl.go`
```go
package services

import (
    "context"
    "fmt"
    "time"
    "gofiber-smart-trash/domain/dto"
    "gofiber-smart-trash/domain/models"
    "gofiber-smart-trash/domain/ports"
    "gofiber-smart-trash/domain/repositories"
)

type trashServiceImpl struct {
    trashRepo      repositories.TrashRepository
    storageAdapter ports.StorageAdapter
}

func NewTrashService(trashRepo repositories.TrashRepository, storageAdapter ports.StorageAdapter) *trashServiceImpl {
    return &trashServiceImpl{
        trashRepo:      trashRepo,
        storageAdapter: storageAdapter,
    }
}

func (s *trashServiceImpl) GenerateUploadURL(ctx context.Context, deviceID string) (*dto.UploadURLResponse, error) {
    // Generate unique key: trash/{device_id}/{timestamp}.jpg
    timestamp := time.Now().UnixMilli()
    key := fmt.Sprintf("trash/%s/%d.jpg", deviceID, timestamp)

    // Generate presigned URL using adapter
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

func (s *trashServiceImpl) CreateTrashRecord(ctx context.Context, req *dto.CreateTrashRequest) (*dto.TrashResponse, error) {
    trash := &models.TrashRecord{
        DeviceID:  req.DeviceID,
        ImageURL:  req.ImageURL,
        Latitude:  req.Latitude,
        Longitude: req.Longitude,
    }

    if err := s.trashRepo.Create(ctx, trash); err != nil {
        return nil, fmt.Errorf("failed to create trash record: %w", err)
    }

    return &dto.TrashResponse{
        ID:        trash.ID,
        DeviceID:  trash.DeviceID,
        ImageURL:  trash.ImageURL,
        Latitude:  trash.Latitude,
        Longitude: trash.Longitude,
        CreatedAt: trash.CreatedAt,
    }, nil
}

func (s *trashServiceImpl) GetTrashByID(ctx context.Context, id uint) (*dto.TrashResponse, error) {
    trash, err := s.trashRepo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get trash record: %w", err)
    }

    return &dto.TrashResponse{
        ID:        trash.ID,
        DeviceID:  trash.DeviceID,
        ImageURL:  trash.ImageURL,
        Latitude:  trash.Latitude,
        Longitude: trash.Longitude,
        CreatedAt: trash.CreatedAt,
    }, nil
}

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
            ID:        trash.ID,
            DeviceID:  trash.DeviceID,
            ImageURL:  trash.ImageURL,
            Latitude:  trash.Latitude,
            Longitude: trash.Longitude,
            CreatedAt: trash.CreatedAt,
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
```

### 3. Infrastructure Layer

#### `infrastructure/postgres/trash_repository_impl.go`
```go
package postgres

import (
    "context"
    "gofiber-smart-trash/domain/models"
    "gofiber-smart-trash/domain/repositories"
    "gorm.io/gorm"
)

type trashRepositoryImpl struct {
    db *gorm.DB
}

func NewTrashRepository(db *gorm.DB) repositories.TrashRepository {
    return &trashRepositoryImpl{db: db}
}

func (r *trashRepositoryImpl) Create(ctx context.Context, trash *models.TrashRecord) error {
    return r.db.WithContext(ctx).Create(trash).Error
}

func (r *trashRepositoryImpl) FindByID(ctx context.Context, id uint) (*models.TrashRecord, error) {
    var trash models.TrashRecord
    if err := r.db.WithContext(ctx).First(&trash, id).Error; err != nil {
        return nil, err
    }
    return &trash, nil
}

func (r *trashRepositoryImpl) FindAll(ctx context.Context, filter repositories.TrashFilter) ([]models.TrashRecord, int64, error) {
    var trashList []models.TrashRecord
    var total int64

    query := r.db.WithContext(ctx).Model(&models.TrashRecord{})

    // Apply filters
    if filter.DeviceID != "" {
        query = query.Where("device_id = ?", filter.DeviceID)
    }

    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply pagination
    if err := query.
        Order("created_at DESC").
        Limit(filter.Limit).
        Offset(filter.Offset).
        Find(&trashList).Error; err != nil {
        return nil, 0, err
    }

    return trashList, total, nil
}
```

#### `infrastructure/storage/r2_adapter.go`
```go
package storage

import (
    "context"
    "fmt"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "gofiber-smart-trash/domain/ports"
)

type R2StorageAdapter struct {
    client    *s3.Client
    bucket    string
    publicURL string
}

func NewR2StorageAdapter(accountID, accessKeyID, secretAccessKey, bucket, publicURL string) (ports.StorageAdapter, error) {
    endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

    cfg, err := config.LoadDefaultConfig(context.Background(),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            accessKeyID,
            secretAccessKey,
            "",
        )),
        config.WithRegion("auto"),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to load AWS config: %w", err)
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.BaseEndpoint = aws.String(endpoint)
    })

    return &R2StorageAdapter{
        client:    client,
        bucket:    bucket,
        publicURL: publicURL,
    }, nil
}

func (a *R2StorageAdapter) GeneratePresignedUploadURL(ctx context.Context, key string, expiry time.Duration) (*ports.PresignedURLResponse, error) {
    presignClient := s3.NewPresignClient(a.client)

    req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(a.bucket),
        Key:         aws.String(key),
        ContentType: aws.String("image/jpeg"),
    }, s3.WithPresignExpires(expiry))

    if err != nil {
        return nil, fmt.Errorf("failed to presign PutObject: %w", err)
    }

    publicURL := a.GeneratePublicURL(key)

    return &ports.PresignedURLResponse{
        UploadURL: req.URL,
        PublicURL: publicURL,
        ExpiresIn: int64(expiry.Seconds()),
    }, nil
}

func (a *R2StorageAdapter) GeneratePublicURL(key string) string {
    return fmt.Sprintf("%s/%s", a.publicURL, key)
}

func (a *R2StorageAdapter) DeleteObject(ctx context.Context, key string) error {
    _, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket: aws.String(a.bucket),
        Key:    aws.String(key),
    })
    return err
}
```

### 4. Interface Layer

#### `interfaces/api/handlers/upload_handler.go`
```go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "gofiber-smart-trash/domain/dto"
)

func (h *Handlers) GenerateUploadURL(c *fiber.Ctx) error {
    deviceID := c.Query("device_id")
    if deviceID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
            Success: false,
            Error:   "MISSING_DEVICE_ID",
            Message: "device_id is required",
        })
    }

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
```

#### `interfaces/api/handlers/trash_handler.go`
```go
package handlers

import (
    "strconv"

    "github.com/gofiber/fiber/v2"
    "gofiber-smart-trash/domain/dto"
    "gofiber-smart-trash/pkg/utils"
)

func (h *Handlers) CreateTrash(c *fiber.Ctx) error {
    var req dto.CreateTrashRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
            Success: false,
            Error:   "INVALID_REQUEST",
            Message: err.Error(),
        })
    }

    if err := utils.ValidateStruct(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
            Success: false,
            Error:   "VALIDATION_ERROR",
            Message: err.Error(),
        })
    }

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

func (h *Handlers) GetTrash(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
            Success: false,
            Error:   "INVALID_ID",
            Message: "Invalid trash ID",
        })
    }

    response, err := h.trashService.GetTrashByID(c.Context(), uint(id))
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

func (h *Handlers) ListTrash(c *fiber.Ctx) error {
    var req dto.ListTrashRequest
    if err := c.QueryParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
            Success: false,
            Error:   "INVALID_REQUEST",
            Message: err.Error(),
        })
    }

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
        Data:    response.Data,
        // Note: Pagination is embedded in ListTrashResponse
    })
}
```

#### `interfaces/api/handlers/handlers.go`
```go
package handlers

import (
    "gofiber-smart-trash/domain/services"
)

type Handlers struct {
    trashService services.TrashService
}

func NewHandlers(trashService services.TrashService) *Handlers {
    return &Handlers{
        trashService: trashService,
    }
}
```

#### `interfaces/api/routes/routes.go`
```go
package routes

import (
    "github.com/gofiber/fiber/v2"
    "gofiber-smart-trash/interfaces/api/handlers"
    "gofiber-smart-trash/interfaces/api/middleware"
)

func SetupRoutes(app *fiber.App, h *handlers.Handlers) {
    // Middleware
    app.Use(middleware.Logger())
    app.Use(middleware.CORS())

    // Health check
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Smart Trash Picker API",
            "version": "1.0.0",
        })
    })

    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "ok",
        })
    })

    // API v1
    api := app.Group("/api")

    // Upload URL generation
    api.Get("/upload-url", h.GenerateUploadURL)

    // Trash management
    api.Post("/trash", h.CreateTrash)
    api.Get("/trash", h.ListTrash)
    api.Get("/trash/:id", h.GetTrash)
}
```

---

## 🔧 Files to MODIFY

### 1. `cmd/api/main.go`

แก้ไขให้ใช้โครงสร้างใหม่:
- เอา Redis, WebSocket, Scheduler ออก
- เพิ่ม Storage Adapter configuration
- Simplify DI container

### 2. `pkg/config/config.go`

ปรับ config structure:
- เอา JWT, Redis config ออก
- เพิ่ม StorageConfig (รองรับหลาย provider)

### 3. `pkg/di/container.go`

Simplify DI:
- เหลือเฉพาะ: DB, Storage Adapter, Repository, Service, Handler
- เอา User, Task, Job, WebSocket, Scheduler ออก

### 4. `infrastructure/postgres/database.go`

ปรับ migration:
- เหลือเฉพาะ `models.TrashRecord`

### 5. `.env.example`

ปรับ environment variables:
```env
# Server
PORT=3000
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=smartpicker
DB_SSL_MODE=disable

# Storage
STORAGE_PROVIDER=r2          # r2, s3, gcs
STORAGE_BUCKET=smart-picker-bucket
STORAGE_PUBLIC_URL=https://pub-xxx.r2.dev

# Cloudflare R2
R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=your_access_key_id
R2_SECRET_ACCESS_KEY=your_secret_access_key

# AWS S3 (if using S3)
# AWS_REGION=us-east-1
# AWS_ACCESS_KEY_ID=your_access_key
# AWS_SECRET_ACCESS_KEY=your_secret_key

# GCS (if using GCS)
# GCS_PROJECT_ID=your_project_id
# GCS_CREDENTIALS_PATH=/path/to/credentials.json
```

---

## 📝 Implementation Steps

### Phase 1: Clean Up (ลบส่วนที่ไม่ใช้)
1. ✅ ลบไฟล์ User, Task, Job, File related
2. ✅ ลบ Redis, WebSocket, Scheduler
3. ✅ ลบ middleware ที่ไม่ใช้ (auth, error)
4. ✅ ลบ routes ที่ไม่ใช้

### Phase 2: Create New Structure
1. ✅ สร้าง domain layer (models, repositories, services, ports, dto)
2. ✅ สร้าง application layer (service implementation)
3. ✅ สร้าง infrastructure layer (repository, R2 adapter)
4. ✅ สร้าง interfaces layer (handlers, routes)

### Phase 3: Update Existing Files
1. ✅ แก้ main.go
2. ✅ แก้ config.go
3. ✅ แก้ container.go
4. ✅ แก้ database.go
5. ✅ แก้ .env.example

### Phase 4: Testing
1. ✅ Test database connection
2. ✅ Test R2 adapter
3. ✅ Test API endpoints
4. ✅ Test end-to-end flow

---

## ✅ Benefits of This Refactoring

### 1. **Simplified Architecture**
- เหลือเฉพาะส่วนที่จำเป็น
- ง่ายต่อการเข้าใจและบำรุงรักษา
- ลด cognitive load

### 2. **Adapter Pattern for Storage**
- เปลี่ยน storage provider ได้ง่าย
- ไม่ต้องแก้ business logic
- รองรับ R2, S3, GCS ในอนาคต

### 3. **Clean Architecture Principles**
- Separation of concerns
- Dependency inversion
- Testability

### 4. **Flexibility**
```go
// ตอนนี้ใช้ R2
adapter := NewR2StorageAdapter(...)

// อนาคตเปลี่ยนเป็น S3 (แค่เปลี่ยน 1 บรรทัด!)
adapter := NewS3StorageAdapter(...)

// หรือ GCS
adapter := NewGCSStorageAdapter(...)

// Business logic ไม่ต้องเปลี่ยนเลย!
service := NewTrashService(repo, adapter)
```

---

## 🎯 Summary

| Metric | Before | After |
|--------|--------|-------|
| Models | 4 (User, Task, File, Job) | 1 (TrashRecord) |
| Services | 4 | 1 (TrashService) |
| Repositories | 4 | 1 (TrashRepository) |
| Handlers | 5 | 2 (Upload, Trash) |
| API Endpoints | 15+ | 4 |
| Infrastructure | Redis, WebSocket, Scheduler, BunnyStorage | PostgreSQL, Storage Adapter |
| LOC | ~5000+ | ~2000 (ประมาณ) |

**Result**: โครงสร้างที่เรียบง่าย ตรงจุด และยืดหยุ่น พร้อมรองรับการเปลี่ยน storage provider ในอนาคต

---

## 🚀 Ready to Start?

เมื่อ approve แผนนี้แล้ว ผมจะเริ่ม:

1. ลบไฟล์ที่ไม่ใช้
2. สร้างไฟล์ใหม่ตามโครงสร้าง
3. ปรับไฟล์เดิมที่ต้องแก้ไข
4. Test ทุก endpoint

พร้อมเริ่มไหมครับ? 🚀
