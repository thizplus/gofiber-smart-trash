-- Migration: Change trash_records.id from BIGSERIAL to UUID
-- Date: 2025-12-04
-- Description: Migrate from auto-increment integer ID to UUID for better security and scalability

-- Step 1: Drop existing table if you want to start fresh (WARNING: This will delete all data!)
DROP TABLE IF EXISTS trash_records CASCADE;

-- Step 2: Create table with UUID primary key
CREATE TABLE trash_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id VARCHAR(20) NOT NULL,
    image_url TEXT NOT NULL,
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Step 3: Create indexes
CREATE INDEX idx_trash_records_device_id ON trash_records(device_id);
CREATE INDEX idx_trash_records_deleted_at ON trash_records(deleted_at);
CREATE INDEX idx_trash_records_created_at ON trash_records(created_at DESC);

-- Step 4: Add comments
COMMENT ON TABLE trash_records IS 'Stores trash collection records from ESP32-CAM devices';
COMMENT ON COLUMN trash_records.id IS 'UUID primary key';
COMMENT ON COLUMN trash_records.device_id IS 'ESP32 device identifier (MAC address)';
COMMENT ON COLUMN trash_records.image_url IS 'Public URL of trash image in R2';
COMMENT ON COLUMN trash_records.latitude IS 'GPS latitude coordinate';
COMMENT ON COLUMN trash_records.longitude IS 'GPS longitude coordinate';
COMMENT ON COLUMN trash_records.created_at IS 'Record creation timestamp';
