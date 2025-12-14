# Smart Trash Picker API Documentation

## Overview

REST API สำหรับระบบรวบรวมข้อมูลขยะ ใช้ GoFiber framework พร้อม PostgreSQL และ Cloudflare R2 storage

**Base URL**: `http://localhost:8080`

---

## API Endpoints

### 1. Root & Health Check

#### GET /
แสดงข้อมูลพื้นฐานของ API

**Response** (200 OK):
```json
{
  "message": "Smart Trash Picker API",
  "version": "1.0.0"
}
```

---

#### GET /health
ตรวจสอบสถานะของ API

**Response** (200 OK):
```json
{
  "status": "ok"
}
```

---

### 2. Upload API

#### GET /api/upload-url
สร้าง Presigned URL สำหรับ upload รูปขยะไปที่ Cloudflare R2

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| device_id | string | Yes | รหัสอุปกรณ์ |

**Request**:
```
GET /api/upload-url?device_id=DEVICE001
```

**Response สำเร็จ** (200 OK):
```json
{
  "success": true,
  "data": {
    "upload_url": "https://xxx.r2.cloudflarestorage.com/bucket/trash/DEVICE001/1702468800000.jpg?X-Amz-...",
    "image_url": "https://pub-xxx.r2.dev/trash/DEVICE001/1702468800000.jpg",
    "expires_in": 900
  }
}
```

**Response ผิดพลาด** (400 Bad Request):
```json
{
  "success": false,
  "error": "MISSING_DEVICE_ID",
  "message": "device_id is required"
}
```

**วิธีใช้งาน Upload URL**:
```bash
# ใช้ upload_url ทำ PUT request upload รูป
curl -X PUT \
  -H "Content-Type: image/jpeg" \
  --data-binary @photo.jpg \
  "https://xxx.r2.cloudflarestorage.com/bucket/trash/DEVICE001/1702468800000.jpg?X-Amz-..."
```

---

### 3. Trash Records API

#### POST /api/trash
สร้างบันทึกข้อมูลขยะใหม่

**Request Body**:
```json
{
  "device_id": "DEVICE001",
  "image_url": "https://pub-xxx.r2.dev/trash/DEVICE001/1702468800000.jpg",
  "latitude": 13.736717,
  "longitude": 100.523186
}
```

**Validation Rules**:
| Field | Type | Required | Validation |
|-------|------|----------|------------|
| device_id | string | Yes | - |
| image_url | string | Yes | ต้องเป็น URL format |
| latitude | float64 | Yes | - |
| longitude | float64 | Yes | - |

**Response สำเร็จ** (201 Created):
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "device_id": "DEVICE001",
    "image_url": "https://pub-xxx.r2.dev/trash/DEVICE001/1702468800000.jpg",
    "latitude": 13.736717,
    "longitude": 100.523186,
    "created_at": "2025-12-13T12:00:00Z"
  }
}
```

**Response ผิดพลาด** (400 Bad Request):
```json
{
  "success": false,
  "error": "VALIDATION_ERROR",
  "message": "image_url: must be a valid URL"
}
```

---

#### GET /api/trash
ดึงรายการข้อมูลขยะทั้งหมด (รองรับ filter และ pagination)

**Query Parameters**:
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| device_id | string | No | - | กรองตาม device_id |
| limit | int | No | 20 | จำนวนรายการต่อหน้า (max: 100) |
| offset | int | No | 0 | ข้ามรายการ |

**Request**:
```
GET /api/trash?device_id=DEVICE001&limit=10&offset=0
```

**Response สำเร็จ** (200 OK):
```json
{
  "success": true,
  "data": {
    "data": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "device_id": "DEVICE001",
        "image_url": "https://pub-xxx.r2.dev/trash/DEVICE001/1702468800000.jpg",
        "latitude": 13.736717,
        "longitude": 100.523186,
        "created_at": "2025-12-13T12:00:00Z"
      },
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "device_id": "DEVICE001",
        "image_url": "https://pub-xxx.r2.dev/trash/DEVICE001/1702468800001.jpg",
        "latitude": 13.736718,
        "longitude": 100.523187,
        "created_at": "2025-12-13T11:00:00Z"
      }
    ],
    "pagination": {
      "total": 100,
      "limit": 10,
      "offset": 0
    }
  }
}
```

---

#### GET /api/trash/:id
ดึงข้อมูลขยะตาม ID

**URL Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | UUID | Yes | UUID ของบันทึกขยะ |

**Request**:
```
GET /api/trash/550e8400-e29b-41d4-a716-446655440000
```

**Response สำเร็จ** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "device_id": "DEVICE001",
    "image_url": "https://pub-xxx.r2.dev/trash/DEVICE001/1702468800000.jpg",
    "latitude": 13.736717,
    "longitude": 100.523186,
    "created_at": "2025-12-13T12:00:00Z"
  }
}
```

**Response ผิดพลาด** (400 Bad Request - Invalid UUID):
```json
{
  "success": false,
  "error": "INVALID_ID",
  "message": "Invalid UUID format"
}
```

**Response ผิดพลาด** (404 Not Found):
```json
{
  "success": false,
  "error": "NOT_FOUND",
  "message": "Trash record not found"
}
```

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| INVALID_REQUEST | 400 | Body parser error |
| VALIDATION_ERROR | 400 | Validation failed |
| INVALID_ID | 400 | UUID format error |
| MISSING_DEVICE_ID | 400 | Required parameter missing |
| NOT_FOUND | 404 | Resource not found |
| INTERNAL_ERROR | 500 | Server error |

---

## Flow การใช้งาน

### 1. Upload รูปภาพและบันทึกข้อมูลขยะ

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. ขอ Presigned URL                                             │
│    GET /api/upload-url?device_id=DEVICE001                      │
│    ← รับ upload_url และ image_url                               │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 2. Upload รูปไป Cloudflare R2                                   │
│    PUT {upload_url}                                             │
│    Content-Type: image/jpeg                                     │
│    Body: <binary image data>                                    │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 3. สร้างบันทึกขยะ                                                │
│    POST /api/trash                                              │
│    Body: { device_id, image_url, latitude, longitude }          │
│    ← รับข้อมูลที่สร้างพร้อม ID                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 2. ดึงข้อมูลขยะ

```
┌─────────────────────────────────────────────────────────────────┐
│ ดึงรายการทั้งหมด                                                 │
│    GET /api/trash                                               │
│    GET /api/trash?device_id=DEVICE001                           │
│    GET /api/trash?limit=10&offset=20                            │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ ดึงตาม ID                                                       │
│    GET /api/trash/550e8400-e29b-41d4-a716-446655440000          │
└─────────────────────────────────────────────────────────────────┘
```

---

## โครงสร้างโปรเจ็ค

```
gofiber-smart-trash/
├── cmd/api/
│   └── main.go                  # Entry point
├── domain/
│   ├── models/                  # Database models
│   │   └── trash.go
│   ├── dto/                     # Request/Response DTOs
│   │   ├── common.go
│   │   └── trash.go
│   ├── services/                # Service interfaces
│   │   └── trash_service.go
│   ├── repositories/            # Repository interfaces
│   │   └── trash_repository.go
│   └── ports/                   # Adapter interfaces
│       └── storage_adapter.go
├── application/
│   └── services/                # Service implementations
│       └── trash_service_impl.go
├── infrastructure/
│   ├── postgres/                # Database implementation
│   │   ├── database.go
│   │   └── trash_repository_impl.go
│   └── storage/                 # Cloud storage
│       └── r2_adapter.go
├── interfaces/api/
│   ├── handlers/                # HTTP handlers
│   │   ├── handlers.go
│   │   ├── trash_handler.go
│   │   └── upload_handler.go
│   ├── middleware/              # Middleware
│   │   ├── cors_middleware.go
│   │   └── logger_middleware.go
│   └── routes/                  # Routes
│       └── routes.go
└── pkg/
    ├── config/                  # Configuration
    │   └── config.go
    ├── di/                      # Dependency Injection
    │   └── container.go
    └── utils/                   # Utilities
        ├── response.go
        └── validator.go
```

---

## Database Schema

**Table: trash_records**

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | รหัสบันทึก |
| device_id | VARCHAR(20) | NOT NULL, INDEX | รหัสอุปกรณ์ |
| image_url | TEXT | NOT NULL | URL รูปภาพ |
| latitude | DECIMAL(10,8) | NOT NULL | พิกัด latitude |
| longitude | DECIMAL(11,8) | NOT NULL | พิกัด longitude |
| created_at | TIMESTAMP | NOT NULL | เวลาสร้าง |
| updated_at | TIMESTAMP | NOT NULL | เวลาแก้ไข |
| deleted_at | TIMESTAMP | NULLABLE, INDEX | Soft delete |

---

## Environment Variables

```env
# Application
APP_NAME=Smart Trash Picker API
PORT=8080
ENV=development

# Database (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=smartpicker
DB_SSL_MODE=disable

# Storage (Cloudflare R2)
STORAGE_PROVIDER=r2
R2_ACCOUNT_ID=xxx
R2_ACCESS_KEY_ID=xxx
R2_SECRET_ACCESS_KEY=xxx
R2_BUCKET=suekk-bucket
R2_PUBLIC_URL=https://pub-xxx.r2.dev
PRESIGNED_URL_EXPIRY=900
```

---

## Technology Stack

- **Framework**: GoFiber v2.52.0
- **Language**: Go 1.23
- **Database**: PostgreSQL + GORM v1.25.6
- **Storage**: Cloudflare R2 (S3-compatible)
- **Validation**: go-playground/validator v10
