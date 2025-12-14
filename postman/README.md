# 📮 Postman Collection - Smart Trash Picker API

## 📁 Files

- **Smart-Trash-Picker-API.postman_collection.json** - API Collection
- **Smart-Trash-Local.postman_environment.json** - Local Environment

## 🚀 Quick Start

### 1. Import to Postman

**Import Collection & Environment:**
1. Open Postman
2. Click **Import** button
3. Drag & drop both files:
   - `Smart-Trash-Picker-API.postman_collection.json`
   - `Smart-Trash-Local.postman_environment.json`

### 2. Select Environment

1. In Postman, select **Smart Trash - Local** from environment dropdown (top-right)
2. Verify variables:
   - `base_url`: `http://localhost:3000`
   - `device_id`: `AABBCC112233`

### 3. Start Server

```bash
# Run server
go run cmd/api/main.go

# Or use Air for hot reload
air

# Or use executable
./smart-trash-api.exe
```

### 4. Test Endpoints

Run requests in this order:

#### Health Check
1. **Root - API Info** → Get API info
2. **Health Check** → Verify server is running

#### Complete ESP32 Flow
Navigate to **"ESP32 Complete Flow"** folder and run in sequence:

1. **1. Generate Upload URL**
   - Automatically saves `upload_url` and `image_url` to environment

2. **2. Upload Image to R2 (Manual)**
   - Body → Binary → Select an image file
   - Click Send

3. **3. Create Trash Record**
   - Automatically uses `image_url` from step 1
   - GPS coordinates are in the body (can be changed)
   - Saves `trash_id` to environment

#### View Records
4. **List All Trash Records** → See all records
5. **List Trash by Device ID** → Filter by device
6. **Get Trash Record by ID** → Get specific record

---

## 📋 API Endpoints

### Health & Info
```
GET  /                    # API Info
GET  /health              # Health Check
```

### Upload Flow
```
GET  /api/upload-url      # Generate Presigned URL
PUT  {upload_url}         # Upload to R2 (external)
```

### Trash Records
```
POST /api/trash           # Create Record
GET  /api/trash           # List Records
GET  /api/trash/:id       # Get Record
```

---

## 🔄 Complete ESP32 Flow

This simulates what the ESP32-CAM device does:

```
┌─────────────────────────────────────────────────────────┐
│ 1. ESP32: GET /api/upload-url?device_id=AABBCC112233   │
│    Response: { upload_url, image_url, expires_in }     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│ 2. ESP32: PUT {upload_url}                              │
│    Body: [JPEG binary data]                             │
│    → Upload directly to Cloudflare R2                   │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│ 3. ESP32: POST /api/trash                               │
│    Body: {                                              │
│      device_id: "AABBCC112233",                         │
│      image_url: "https://pub-xxx.r2.dev/...",          │
│      latitude: 13.756331,                               │
│      longitude: 100.501762                              │
│    }                                                    │
└─────────────────────────────────────────────────────────┘
```

---

## 🧪 Testing Tips

### Test with Different Device IDs

Change `device_id` in environment:
```
AABBCC112233 (default)
DDEEFF445566
ESP32CAM001
```

### Test Pagination

Modify query parameters:
```
limit: 10, 20, 50, 100
offset: 0, 10, 20, ...
```

### Test GPS Coordinates

Different locations in Bangkok:
```json
// Siam
{ "latitude": 13.745686, "longitude": 100.534821 }

// Chatuchak
{ "latitude": 13.798863, "longitude": 100.551148 }

// Sukhumvit
{ "latitude": 13.736717, "longitude": 100.560471 }
```

### Upload Different Images

In "Upload Image to R2" request:
1. Body → binary
2. Select file → Choose different JPEG files
3. Send

---

## 📝 Environment Variables

### Collection Variables (Shared)
```
base_url       = http://localhost:3000
device_id      = AABBCC112233
limit          = 20
offset         = 0
```

### Auto-Saved Variables (Scripts)
```
upload_url     = (saved from Generate Upload URL)
image_url      = (saved from Generate Upload URL)
expires_in     = (saved from Generate Upload URL)
trash_id       = (saved from Create Trash Record)
```

---

## 🔧 Troubleshooting

### Error: "No upload_url found"
→ Run **"Generate Upload URL"** first

### Error: "Connection refused"
→ Make sure server is running on port 3000

### Error: "MISSING_DEVICE_ID"
→ Check if `device_id` query parameter is set

### Error: "VALIDATION_ERROR"
→ Verify request body format matches schema

### Upload to R2 fails
→ Check if presigned URL expired (15 min limit)
→ Re-run "Generate Upload URL"

---

## 📚 Additional Resources

- **API Documentation**: See REFACTORING_PLAN.md
- **Architecture**: Clean Architecture + Adapter Pattern
- **Database Schema**: See smart-trash-picker-plan-v3.md

---

## 🎉 Happy Testing!
