# Smart Trash API - Test Results

ทดสอบเมื่อ: **2025-12-13 20:24:36 (UTC+7)**

---

## Step 1: ขอ Presigned URL

### Request
```http
GET /api/upload-url?device_id=TEST001 HTTP/1.1
Host: localhost:8080
```

### curl Command
```bash
curl -X GET "http://localhost:8080/api/upload-url?device_id=TEST001"
```

### Response
**Status**: `200 OK`

```json
{
  "success": true,
  "data": {
    "upload_url": "https://suekk-bucket.fcc0e164ed5f9fcf121a73f8f111ccd1.r2.cloudflarestorage.com/trash/TEST001/1765632276936.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ec1293bc10ab6bff0c09788a73b839dd%2F20251213%2Fauto%2Fs3%2Faws4_request&X-Amz-Date=20251213T132436Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host&x-id=PutObject&X-Amz-Signature=f03f55d456b7b2c271cb5f7f0e1961d28e93266da0817c3a14c0ef543cbabe41",
    "image_url": "https://pub-a058b390b77f486aaf97a1d1f073c6c8.r2.dev/trash/TEST001/1765632276936.jpg",
    "expires_in": 900
  }
}
```

### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| success | boolean | สถานะการทำงาน |
| data.upload_url | string | Presigned URL สำหรับ PUT upload รูป (หมดอายุ 15 นาที) |
| data.image_url | string | Public URL ของรูปภาพ (ใช้บันทึกลง database) |
| data.expires_in | number | เวลาหมดอายุของ upload_url (วินาที) |

---

## Step 2: Upload รูปไป Cloudflare R2

### Request
```http
PUT /trash/TEST001/1765632276936.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&... HTTP/1.1
Host: suekk-bucket.fcc0e164ed5f9fcf121a73f8f111ccd1.r2.cloudflarestorage.com
Content-Type: image/jpeg
Content-Length: 399274

<binary image data>
```

### curl Command
```bash
curl -X PUT \
  -H "Content-Type: image/jpeg" \
  --data-binary "@xxx.jpg" \
  "https://suekk-bucket.fcc0e164ed5f9fcf121a73f8f111ccd1.r2.cloudflarestorage.com/trash/TEST001/1765632276936.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ec1293bc10ab6bff0c09788a73b839dd%2F20251213%2Fauto%2Fs3%2Faws4_request&X-Amz-Date=20251213T132436Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host&x-id=PutObject&X-Amz-Signature=f03f55d456b7b2c271cb5f7f0e1961d28e93266da0817c3a14c0ef543cbabe41"
```

### Response
**Status**: `200 OK`

### Response Headers from R2
```
HTTP/1.1 200 OK
Date: Sat, 13 Dec 2025 13:25:00 GMT
Content-Type: text/plain;charset=UTF-8
Content-Length: 0
Connection: keep-alive
ETag: "bc8ba6fa4768d05155bbf19e3485b452"
x-amz-checksum-crc64nvme: inWgSj9VgCs=
x-amz-version-id: 7e64e81d2d79a0a9d943b6fea3f3fd7d
Server: cloudflare
CF-RAY: 9ad5c8ec7e15d02b-BKK
```

### Response Headers Description
| Header | Value | Description |
|--------|-------|-------------|
| HTTP Status | 200 OK | Upload สำเร็จ |
| ETag | "bc8ba6fa4768d05155bbf19e3485b452" | MD5 hash ของไฟล์ที่ upload |
| x-amz-version-id | 7e64e81d2d79a0a9d943b6fea3f3fd7d | Version ID ของไฟล์ใน R2 |
| x-amz-checksum-crc64nvme | inWgSj9VgCs= | CRC64 checksum |
| Server | cloudflare | Server ที่ให้บริการ |
| CF-RAY | 9ad5c8ec7e15d02b-BKK | Cloudflare Ray ID (BKK = Bangkok) |

---

## Step 3: สร้างบันทึกขยะ

### Request
```http
POST /api/trash HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "device_id": "TEST001",
  "image_url": "https://pub-a058b390b77f486aaf97a1d1f073c6c8.r2.dev/trash/TEST001/1765632276936.jpg",
  "latitude": 13.736717,
  "longitude": 100.523186
}
```

### curl Command
```bash
curl -X POST "http://localhost:8080/api/trash" \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "TEST001",
    "image_url": "https://pub-a058b390b77f486aaf97a1d1f073c6c8.r2.dev/trash/TEST001/1765632276936.jpg",
    "latitude": 13.736717,
    "longitude": 100.523186
  }'
```

### Request Body Fields
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| device_id | string | Yes | รหัสอุปกรณ์ |
| image_url | string | Yes | URL รูปภาพ (จาก Step 1) |
| latitude | number | Yes | พิกัด latitude |
| longitude | number | Yes | พิกัด longitude |

### Response
**Status**: `201 Created`

### Response Headers
```
HTTP/1.1 201 Created
Date: Sat, 13 Dec 2025 13:25:27 GMT
Content-Type: application/json
Content-Length: 282
Vary: Origin
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
```

### Response Body
```json
{
  "success": true,
  "data": {
    "id": "51eb676e-9c63-49f3-bad7-186f4bea1e6e",
    "device_id": "TEST001",
    "image_url": "https://pub-a058b390b77f486aaf97a1d1f073c6c8.r2.dev/trash/TEST001/1765632276936.jpg",
    "latitude": 13.736717,
    "longitude": 100.523186,
    "created_at": "2025-12-13T20:25:27.9574523+07:00"
  }
}
```

### Response Body Fields
| Field | Type | Description |
|-------|------|-------------|
| success | boolean | สถานะการทำงาน |
| data.id | string (UUID) | รหัสบันทึกขยะ (auto-generated) |
| data.device_id | string | รหัสอุปกรณ์ |
| data.image_url | string | URL รูปภาพ |
| data.latitude | number | พิกัด latitude |
| data.longitude | number | พิกัด longitude |
| data.created_at | string (ISO 8601) | เวลาที่สร้างบันทึก |

---

## Summary

| Step | Endpoint | Method | Status | Result |
|------|----------|--------|--------|--------|
| 1 | `/api/upload-url` | GET | 200 OK | ได้ presigned URL |
| 2 | R2 Storage | PUT | 200 OK | Upload รูปสำเร็จ |
| 3 | `/api/trash` | POST | 201 Created | สร้างบันทึกสำเร็จ |

---

## ตรวจสอบรูปภาพที่ Upload

รูปภาพที่ upload สามารถเข้าถึงได้ที่:
```
https://pub-a058b390b77f486aaf97a1d1f073c6c8.r2.dev/trash/TEST001/1765632276936.jpg
```

---

## Full Flow Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│ Step 1: ขอ Presigned URL                                                 │
│                                                                          │
│   Client ──GET /api/upload-url?device_id=TEST001──> API Server           │
│          <────────── upload_url + image_url ──────────                   │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│ Step 2: Upload รูปไป R2                                                  │
│                                                                          │
│   Client ──PUT {upload_url} + binary image──> Cloudflare R2              │
│          <────────────── 200 OK + ETag ───────────────                   │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│ Step 3: สร้างบันทึกขยะ                                                    │
│                                                                          │
│   Client ──POST /api/trash + {device_id, image_url, lat, lng}──> API     │
│          <────────────── 201 Created + record data ──────────            │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## Test File Info

| Property | Value |
|----------|-------|
| File Name | xxx.jpg |
| File Size | 399,274 bytes (~390 KB) |
| Content-Type | image/jpeg |
| R2 ETag | bc8ba6fa4768d05155bbf19e3485b452 |
