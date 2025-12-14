# 🗑️ Smart Trash Picker - Project Plan v3

## 📋 Overview

ไม้คีบขยะอัจฉริยะที่ถ่ายภาพขยะอัตโนมัติเมื่อคีบ พร้อมบันทึกพิกัด GPS

### Tech Stack

| Layer | Technology |
|-------|------------|
| Hardware | ESP32-CAM + GPS Neo-6M |
| Backend | Go Fiber |
| Database | PostgreSQL |
| Storage | Cloudflare R2 (Presigned URL) |

---

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Smart Trash Picker                                │
└─────────────────────────────────────────────────────────────────────────────┘

    ESP32-CAM                       Go Fiber                    Cloudflare R2
        │                              │                              │
        │  1. กดคีบ → ถ่ายรูป + GPS    │                              │
        │                              │                              │
        │  2. GET /api/upload-url      │                              │
        │  ───────────────────────────►│                              │
        │                              │                              │
        │                              │  สร้างชื่อไฟล์               │
        │                              │  สร้าง Presigned URL         │
        │                              │                              │
        │  3. Response                 │                              │
        │  {                           │                              │
        │    upload_url: "...",        │                              │
        │    image_url: "..."          │                              │
        │  }                           │                              │
        │  ◄───────────────────────────│                              │
        │                              │                              │
        │  4. PUT ภาพไป upload_url     │                              │
        │  ────────────────────────────────────────────────────────────►
        │                              │                              │
        │                              │                       ภาพถูกเก็บ
        │                              │                              │
        │  5. POST /api/trash          │                              │
        │  {                           │                              │
        │    device_id,                │                              │
        │    image_url,                │                              │
        │    lat, lng                  │                              │
        │  }                           │                              │
        │  ───────────────────────────►│                              │
        │                              │                              │
        │                              │  บันทึก PostgreSQL           │
        │                              │                              │
        │  6. Response                 │                              │
        │  {success: true}             │                              │
        │  ◄───────────────────────────│                              │
        │                              │                              │
```

---

## 🔄 Complete Flow

### Step-by-Step

```
┌─────────────────────────────────────────────────────────────────┐
│ Step 1: กดคีบขยะ                                                 │
├─────────────────────────────────────────────────────────────────┤
│ - Micro Switch ถูกกด                                            │
│ - ESP32 ถ่ายภาพ (JPEG)                                          │
│ - ESP32 อ่าน GPS (latitude, longitude)                          │
│ - เก็บภาพไว้ใน memory                                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 2: ขอ Presigned URL                                        │
├─────────────────────────────────────────────────────────────────┤
│ ESP32 Request:                                                  │
│   GET /api/upload-url?device_id=AABBCC112233                    │
│                                                                 │
│ Backend ทำงาน:                                                  │
│   1. รับ device_id                                              │
│   2. สร้างชื่อไฟล์: trash/{device_id}/{timestamp}.jpg           │
│   3. สร้าง Presigned URL (หมดอายุ 15 นาที)                      │
│   4. สร้าง Public URL                                           │
│                                                                 │
│ Backend Response:                                               │
│   {                                                             │
│     "upload_url": "https://{account}.r2.cloudflarestorage.com   │
│                    /bucket/trash/AABBCC112233/1701590400000.jpg │
│                    ?X-Amz-Algorithm=AWS4-HMAC-SHA256            │
│                    &X-Amz-Credential=...                        │
│                    &X-Amz-Signature=...",                       │
│     "image_url": "https://pub-xxx.r2.dev                        │
│                   /trash/AABBCC112233/1701590400000.jpg",       │
│     "expires_in": 900                                           │
│   }                                                             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 3: Upload ภาพไป R2                                         │
├─────────────────────────────────────────────────────────────────┤
│ ESP32 Request:                                                  │
│   PUT {upload_url}                                              │
│   Content-Type: image/jpeg                                      │
│   Body: [JPEG binary data]                                      │
│                                                                 │
│ R2 Response:                                                    │
│   HTTP 200 OK                                                   │
│                                                                 │
│ ภาพถูกเก็บที่: trash/AABBCC112233/1701590400000.jpg             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 4: บันทึกข้อมูลลง Database                                  │
├─────────────────────────────────────────────────────────────────┤
│ ESP32 Request:                                                  │
│   POST /api/trash                                               │
│   {                                                             │
│     "device_id": "AABBCC112233",                                │
│     "image_url": "https://pub-xxx.r2.dev/trash/.../xxx.jpg",    │
│     "latitude": 13.756331,                                      │
│     "longitude": 100.501762                                     │
│   }                                                             │
│                                                                 │
│ Backend ทำงาน:                                                  │
│   1. Validate ข้อมูล                                            │
│   2. INSERT INTO trash_records                                  │
│                                                                 │
│ Backend Response:                                               │
│   {                                                             │
│     "success": true,                                            │
│     "data": {                                                   │
│       "id": 123,                                                │
│       "device_id": "AABBCC112233",                              │
│       "image_url": "...",                                       │
│       "latitude": 13.756331,                                    │
│       "longitude": 100.501762,                                  │
│       "created_at": "2025-12-03T10:30:00Z"                      │
│     }                                                           │
│   }                                                             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 5: LED Feedback                                            │
├─────────────────────────────────────────────────────────────────┤
│ Success: LED กระพริบ 2 ครั้ง (ช้า)                               │
│ Error: LED กระพริบ 5 ครั้ง (เร็ว)                                │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🗄️ Database Schema

### PostgreSQL

```sql
-- สร้าง Database
CREATE DATABASE smartpicker;

-- สร้างตาราง
CREATE TABLE trash_records (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(20) NOT NULL,
    image_url TEXT NOT NULL,
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Index สำหรับ query
CREATE INDEX idx_trash_device_id ON trash_records(device_id);
CREATE INDEX idx_trash_created_at ON trash_records(created_at DESC);
```

### ตัวอย่างข้อมูล

| id | device_id | image_url | latitude | longitude | created_at |
|----|-----------|-----------|----------|-----------|------------|
| 1 | AABBCC112233 | https://pub-xxx.r2.dev/trash/AABBCC112233/1701590400000.jpg | 13.75633100 | 100.50176200 | 2025-12-03 10:30:00 |
| 2 | AABBCC112233 | https://pub-xxx.r2.dev/trash/AABBCC112233/1701590535000.jpg | 13.75645200 | 100.50182300 | 2025-12-03 10:32:15 |
| 3 | DDEEFF445566 | https://pub-xxx.r2.dev/trash/DDEEFF445566/1701591000000.jpg | 13.76012300 | 100.49876500 | 2025-12-03 10:40:00 |

---

## ☁️ Cloudflare R2 Setup

### Bucket Structure

```
smart-picker-bucket/
└── trash/
    ├── AABBCC112233/                    ← device_id
    │   ├── 1701590400000.jpg            ← timestamp.jpg
    │   ├── 1701590535000.jpg
    │   └── 1701591200000.jpg
    │
    └── DDEEFF445566/                    ← another device
        ├── 1701591000000.jpg
        └── 1701592000000.jpg
```

### URL Format

| Type | URL |
|------|-----|
| Presigned (Upload) | `https://{account_id}.r2.cloudflarestorage.com/{bucket}/trash/{device_id}/{timestamp}.jpg?X-Amz-Signature=...` |
| Public (View) | `https://pub-{hash}.r2.dev/trash/{device_id}/{timestamp}.jpg` |

### R2 Configuration

```
Account ID: {your_account_id}
Bucket Name: smart-picker-bucket
Public Access: Enabled (r2.dev subdomain)
```

### Credentials (เก็บที่ Backend เท่านั้น)

```env
R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=your_access_key
R2_SECRET_ACCESS_KEY=your_secret_key
R2_BUCKET_NAME=smart-picker-bucket
R2_PUBLIC_URL=https://pub-xxx.r2.dev
```

> ⚠️ **สำคัญ**: ESP32 ไม่ต้องรู้ R2 credentials เลย ใช้ Presigned URL จาก Backend

---

## 🔌 API Endpoints

### Base URL
```
https://api.smartpicker.example.com
```

---

### 1. GET /api/upload-url

**ขอ Presigned URL สำหรับ upload ภาพ**

Request:
```
GET /api/upload-url?device_id=AABBCC112233
```

Response (Success):
```json
{
    "success": true,
    "data": {
        "upload_url": "https://{account}.r2.cloudflarestorage.com/smart-picker-bucket/trash/AABBCC112233/1701590400000.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...&X-Amz-Date=...&X-Amz-Expires=900&X-Amz-Signature=...",
        "image_url": "https://pub-xxx.r2.dev/trash/AABBCC112233/1701590400000.jpg",
        "expires_in": 900
    }
}
```

Response (Error):
```json
{
    "success": false,
    "error": "MISSING_DEVICE_ID",
    "message": "device_id is required"
}
```

---

### 2. POST /api/trash

**บันทึกข้อมูลขยะ (หลังจาก upload รูปสำเร็จ)**

Request:
```json
{
    "device_id": "AABBCC112233",
    "image_url": "https://pub-xxx.r2.dev/trash/AABBCC112233/1701590400000.jpg",
    "latitude": 13.756331,
    "longitude": 100.501762
}
```

Response (Success):
```json
{
    "success": true,
    "data": {
        "id": 123,
        "device_id": "AABBCC112233",
        "image_url": "https://pub-xxx.r2.dev/trash/AABBCC112233/1701590400000.jpg",
        "latitude": 13.756331,
        "longitude": 100.501762,
        "created_at": "2025-12-03T10:30:00Z"
    }
}
```

Response (Error):
```json
{
    "success": false,
    "error": "VALIDATION_ERROR",
    "message": "latitude is required"
}
```

---

### 3. GET /api/trash

**ดึงรายการขยะทั้งหมด**

Request:
```
GET /api/trash
GET /api/trash?device_id=AABBCC112233
GET /api/trash?limit=20&offset=0
```

Response:
```json
{
    "success": true,
    "data": [
        {
            "id": 123,
            "device_id": "AABBCC112233",
            "image_url": "https://pub-xxx.r2.dev/trash/AABBCC112233/1701590400000.jpg",
            "latitude": 13.756331,
            "longitude": 100.501762,
            "created_at": "2025-12-03T10:30:00Z"
        },
        {
            "id": 122,
            "device_id": "AABBCC112233",
            "image_url": "https://pub-xxx.r2.dev/trash/AABBCC112233/1701590200000.jpg",
            "latitude": 13.756210,
            "longitude": 100.501650,
            "created_at": "2025-12-03T10:26:40Z"
        }
    ],
    "pagination": {
        "total": 156,
        "limit": 20,
        "offset": 0
    }
}
```

---

### 4. GET /api/trash/:id

**ดึงรายการขยะเดียว**

Request:
```
GET /api/trash/123
```

Response:
```json
{
    "success": true,
    "data": {
        "id": 123,
        "device_id": "AABBCC112233",
        "image_url": "https://pub-xxx.r2.dev/trash/AABBCC112233/1701590400000.jpg",
        "latitude": 13.756331,
        "longitude": 100.501762,
        "created_at": "2025-12-03T10:30:00Z"
    }
}
```

---

## 📁 Project Structure

### Backend (Go Fiber)

```
smart-picker-api/
├── main.go                      # Entry point
├── go.mod
├── go.sum
├── .env                         # Environment variables
│
├── config/
│   └── config.go                # Load environment variables
│
├── database/
│   ├── database.go              # PostgreSQL connection
│   └── migrations/
│       └── 001_create_trash_records.sql
│
├── models/
│   └── trash.go                 # Trash struct
│
├── handlers/
│   ├── upload.go                # GET /api/upload-url
│   └── trash.go                 # POST/GET /api/trash
│
├── services/
│   └── r2.go                    # Cloudflare R2 presigned URL
│
├── routes/
│   └── routes.go                # Route definitions
│
└── utils/
    └── response.go              # JSON response helpers
```

---

## 🔌 ESP32 Firmware

### Pin Configuration

| Component | ESP32-CAM Pin | หมายเหตุ |
|-----------|---------------|----------|
| Camera | Built-in | OV2640 |
| GPS TX | GPIO 16 | → ESP32 RX2 |
| GPS RX | GPIO 17 | ← ESP32 TX2 |
| Micro Switch | GPIO 12 | Pull-up, active LOW |
| Status LED | GPIO 4 | Built-in Flash LED |
| Power Switch | External | ต่อกับ battery |

### Wiring Diagram

```
                    ┌─────────────────────────────────────┐
                    │           ESP32-CAM                 │
                    │                                     │
                    │   [OV2640 Camera - Built-in]        │
                    │                                     │
                    │                           3.3V ─────┼──► GPS VCC
                    │                            GND ─────┼──► GPS GND
  Micro Switch ────►│ GPIO 12                            │
       │            │                         GPIO 16 ◄──┼─── GPS TX
       └──► GND     │                         GPIO 17 ───┼──► GPS RX
                    │                                     │
                    │   GPIO 4 = Status LED (built-in)   │
                    │                                     │
                    │                              5V ◄───┼─── Power (MT3608)
                    │                             GND ◄───┼─── Power GND
                    └─────────────────────────────────────┘

Power Circuit:
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌───────────┐
│ 18650   │───►│ TP4056  │───►│ MT3608  │───►│ ESP32-CAM │
│ 3.7V    │    │ Charger │    │ 3.7→5V  │    │           │
└─────────┘    └─────────┘    └─────────┘    └───────────┘
                    │
                    │
              ┌─────────────┐
              │ Power Switch │
              │ (On/Off)     │
              └─────────────┘
```

### Firmware Flow

```
┌─────────────────────────────────────────────┐
│                   START                      │
└─────────────────┬───────────────────────────┘
                  ▼
┌─────────────────────────────────────────────┐
│              Initialize                      │
│  - Serial (debug)                           │
│  - Camera                                   │
│  - GPS                                      │
│  - Switch pin (INPUT_PULLUP)                │
└─────────────────┬───────────────────────────┘
                  ▼
┌─────────────────────────────────────────────┐
│         Connect WiFi (Hotspot)              │
│         Retry until connected               │
│         LED: Blink slow                     │
└─────────────────┬───────────────────────────┘
                  ▼
┌─────────────────────────────────────────────┐
│         Wait for GPS Fix                    │
│         (อาจใช้เวลา 30-60 วินาที)            │
│         LED: Blink medium                   │
└─────────────────┬───────────────────────────┘
                  ▼
┌─────────────────────────────────────────────┐
│              READY                          │
│         LED: Solid ON                       │
└─────────────────┬───────────────────────────┘
                  ▼
┌─────────────────────────────────────────────┐
│    ┌─────────────────────────────────┐      │
│    │   Wait for Switch Press         │◄─────┼───┐
│    └─────────────────────────────────┘      │   │
│                  │ (pressed)                │   │
│                  ▼                          │   │
│    ┌─────────────────────────────────┐      │   │
│    │   LED: Blink fast (working)     │      │   │
│    └─────────────────────────────────┘      │   │
│                  │                          │   │
│                  ▼                          │   │
│    ┌─────────────────────────────────┐      │   │
│    │   1. Capture Photo              │      │   │
│    │   2. Read GPS (lat, lng)        │      │   │
│    └─────────────────────────────────┘      │   │
│                  │                          │   │
│                  ▼                          │   │
│    ┌─────────────────────────────────┐      │   │
│    │   3. GET /api/upload-url        │      │   │
│    │      → upload_url, image_url    │      │   │
│    └─────────────────────────────────┘      │   │
│                  │                          │   │
│                  ▼                          │   │
│    ┌─────────────────────────────────┐      │   │
│    │   4. PUT photo to upload_url    │      │   │
│    └─────────────────────────────────┘      │   │
│                  │                          │   │
│                  ▼                          │   │
│    ┌─────────────────────────────────┐      │   │
│    │   5. POST /api/trash            │      │   │
│    │      {device_id, image_url,     │      │   │
│    │       lat, lng}                 │      │   │
│    └─────────────────────────────────┘      │   │
│                  │                          │   │
│                  ▼                          │   │
│    ┌─────────────────────────────────┐      │   │
│    │   6. LED Feedback               │      │   │
│    │      Success: 2x slow blink     │      │   │
│    │      Error: 5x fast blink       │      │   │
│    └─────────────────────────────────┘      │   │
│                  │                          │   │
│                  └──────────────────────────┼───┘
└─────────────────────────────────────────────┘
```

### Config (config.h)

```cpp
#ifndef CONFIG_H
#define CONFIG_H

// ==================== WiFi ====================
#define WIFI_SSID "YourHotspotName"
#define WIFI_PASSWORD "YourHotspotPassword"
#define WIFI_TIMEOUT 30000  // 30 seconds

// ==================== Backend API ====================
#define API_BASE_URL "https://api.smartpicker.example.com"
#define API_UPLOAD_URL_ENDPOINT "/api/upload-url"
#define API_TRASH_ENDPOINT "/api/trash"

// ==================== Hardware Pins ====================
#define GPS_RX_PIN 16
#define GPS_TX_PIN 17
#define GPS_BAUD 9600

#define SWITCH_PIN 12
#define LED_PIN 4

// ==================== Settings ====================
#define DEBOUNCE_DELAY 200      // ms
#define GPS_TIMEOUT 60000       // 60 seconds
#define HTTP_TIMEOUT 30000      // 30 seconds
#define IMAGE_QUALITY 12        // 0-63, lower = better quality

#endif
```

---

## 📊 Environment Variables

### Backend (.env)

```env
# Server
PORT=3000
ENV=development

# Database
DATABASE_URL=postgres://user:password@localhost:5432/smartpicker

# Cloudflare R2
R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=your_access_key_id
R2_SECRET_ACCESS_KEY=your_secret_access_key
R2_BUCKET_NAME=smart-picker-bucket
R2_PUBLIC_URL=https://pub-xxx.r2.dev

# Presigned URL
PRESIGNED_URL_EXPIRY=900  # 15 minutes in seconds
```

### ESP32 (config.h)

```cpp
// WiFi - เปลี่ยนตามมือถือ
#define WIFI_SSID "MyPhone"
#define WIFI_PASSWORD "12345678"

// Backend - เปลี่ยนตาม server
#define API_BASE_URL "https://api.smartpicker.example.com"
```

> ⚠️ **สังเกต**: ESP32 ไม่มี R2 credentials เลย! ปลอดภัยกว่า

---

## 🚀 Development Steps

### Phase 1: Setup Infrastructure

#### Step 1.1: Create PostgreSQL Database
```bash
# สร้าง database
createdb smartpicker

# หรือใช้ psql
psql -U postgres -c "CREATE DATABASE smartpicker;"
```

#### Step 1.2: Run Migration
```sql
-- File: 001_create_trash_records.sql
CREATE TABLE trash_records (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(20) NOT NULL,
    image_url TEXT NOT NULL,
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_trash_device_id ON trash_records(device_id);
CREATE INDEX idx_trash_created_at ON trash_records(created_at DESC);
```

#### Step 1.3: Setup Cloudflare R2
1. Login Cloudflare Dashboard
2. ไปที่ R2 → Create bucket: `smart-picker-bucket`
3. Settings → Public access → Enable (ใช้ r2.dev subdomain)
4. สร้าง API Token:
   - Permissions: Object Read & Write
   - จด Access Key ID และ Secret Access Key
5. จด Public URL: `https://pub-{hash}.r2.dev`

---

### Phase 2: Backend Development

#### Step 2.1: Initialize Go Project
```bash
mkdir smart-picker-api
cd smart-picker-api
go mod init smart-picker-api
```

#### Step 2.2: Install Dependencies
```bash
go get github.com/gofiber/fiber/v2
go get github.com/joho/godotenv
go get github.com/lib/pq
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/s3
```

#### Step 2.3: Implement APIs
1. GET /api/upload-url - สร้าง presigned URL
2. POST /api/trash - บันทึกข้อมูล
3. GET /api/trash - ดึงรายการ

#### Step 2.4: Test APIs
```bash
# Test upload URL
curl "http://localhost:3000/api/upload-url?device_id=TEST123"

# Test create trash
curl -X POST http://localhost:3000/api/trash \
  -H "Content-Type: application/json" \
  -d '{"device_id":"TEST123","image_url":"https://...","latitude":13.756,"longitude":100.501}'

# Test get all
curl "http://localhost:3000/api/trash"
```

---

### Phase 3: ESP32 Firmware Development

#### Step 3.1: Setup Arduino IDE
1. Install ESP32 board support
2. Select: AI Thinker ESP32-CAM
3. Install libraries:
   - TinyGPSPlus
   - ArduinoJson
   - HTTPClient (built-in)

#### Step 3.2: Test Individual Components
1. Test Camera - ถ่ายรูปและดูใน Serial
2. Test GPS - อ่านพิกัดและดูใน Serial
3. Test WiFi - เชื่อมต่อ hotspot
4. Test HTTP - ส่ง request ไป server

#### Step 3.3: Integrate All Components
1. รวม code ทั้งหมด
2. ทดสอบ flow เต็ม
3. Debug และแก้ไข

#### Step 3.4: Test End-to-End
1. เปิด hotspot มือถือ
2. เปิด ESP32
3. รอ GPS fix
4. กดปุ่ม (Micro Switch)
5. ตรวจสอบ:
   - ภาพใน R2
   - ข้อมูลใน Database

---

### Phase 4: Assembly

#### Step 4.1: Prepare Hardware
1. ต่อวงจรบน breadboard ทดสอบก่อน
2. บัดกรีเมื่อแน่ใจว่าทำงานได้

#### Step 4.2: Mount on Trash Picker
1. ออกแบบเคส (3D print หรือทำเอง)
2. ติดตั้งบนไม้คีบ
3. จัดวางสาย

#### Step 4.3: Final Testing
1. ทดสอบการใช้งานจริง
2. ตรวจสอบความทนทาน
3. ปรับปรุงถ้าจำเป็น

---

## 📝 Files to Create

### Backend

| ไฟล์ | คำอธิบาย |
|------|---------|
| `main.go` | Entry point, setup Fiber |
| `config/config.go` | Load .env |
| `database/database.go` | PostgreSQL connection |
| `models/trash.go` | Trash struct |
| `handlers/upload.go` | GET /api/upload-url |
| `handlers/trash.go` | POST/GET /api/trash |
| `services/r2.go` | R2 presigned URL generator |
| `routes/routes.go` | Route setup |
| `utils/response.go` | JSON response helpers |
| `.env` | Environment variables |

### ESP32 Firmware

| ไฟล์ | คำอธิบาย |
|------|---------|
| `firmware.ino` | Main Arduino sketch |
| `config.h` | Configuration constants |
| `camera.h` | Camera functions |
| `gps.h` | GPS functions |
| `network.h` | WiFi & HTTP functions |

---

## ✅ Checklist

### Hardware
- [ ] ESP32-CAM
- [ ] GPS Neo-6M
- [ ] USB TTL (FT232RL)
- [ ] แบต 18650 3.7V
- [ ] TP4056 (ชาร์จแบต)
- [ ] MT3608 (Step Up 3.7V → 5V)
- [ ] Micro Switch (SS-5GL2)
- [ ] สวิตช์เปิด-ปิด
- [ ] สาย Jumper

### Backend
- [ ] PostgreSQL database created
- [ ] Cloudflare R2 bucket created
- [ ] R2 public access enabled
- [ ] R2 API token created
- [ ] Go Fiber API implemented
- [ ] APIs tested

### ESP32
- [ ] Arduino IDE setup
- [ ] Libraries installed
- [ ] Camera tested
- [ ] GPS tested
- [ ] WiFi tested
- [ ] Full flow tested

### Integration
- [ ] End-to-end test passed
- [ ] Mounted on trash picker
- [ ] Real-world test passed

---

## 🔮 Future Enhancements (Optional)

เมื่อ MVP เสร็จแล้ว สามารถเพิ่มได้:

1. **AI Classification** - วิเคราะห์ประเภทขยะด้วย OpenAI Vision
2. **User System** - ลงทะเบียน device กับ user
3. **Campaign** - แคมเปญแข่งขันเก็บขยะ
4. **Dashboard** - หน้าเว็บแสดงแผนที่และสถิติ
5. **Mobile App** - แอปสำหรับดูข้อมูล

---

## 📚 References

- [ESP32-CAM Documentation](https://randomnerdtutorials.com/esp32-cam-video-streaming-face-recognition-arduino-ide/)
- [TinyGPS++ Library](http://arduiniana.org/libraries/tinygpsplus/)
- [Go Fiber Documentation](https://docs.gofiber.io/)
- [Cloudflare R2 Documentation](https://developers.cloudflare.com/r2/)
- [AWS S3 Presigned URLs](https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-presigned-url.html)

---

## ⏭️ Next Step

เริ่มเขียน code ไหมครับ?

1. **Backend** - Go Fiber + PostgreSQL + R2 Presigned URL
2. **ESP32 Firmware** - Camera + GPS + HTTP

แนะนำเริ่มจาก **Backend** ก่อน เพราะ ESP32 ต้องมี API ให้เรียก
