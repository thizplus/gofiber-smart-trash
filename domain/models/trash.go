package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TrashRecord struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DeviceID  string         `gorm:"type:varchar(20);not null;index" json:"device_id"`
	ImageURL  string         `gorm:"type:text;not null" json:"image_url"`
	Latitude  float64        `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude float64        `gorm:"type:decimal(11,8);not null" json:"longitude"`

	// AI Classification fields
	Category      string    `gorm:"type:varchar(50)" json:"category"`       // cardboard, glass, metal, paper, plastic, trash
	SubCategory   string    `gorm:"type:varchar(50)" json:"sub_category"`   // For L2 classification (e.g., PET, HDPE)
	Confidence    float64   `gorm:"type:decimal(5,4)" json:"confidence"`    // 0.0000 - 1.0000
	BinNumber     int       `gorm:"type:int" json:"bin_number"`             // 1-6
	BinLabel      string    `gorm:"type:varchar(50)" json:"bin_label"`      // Thai label
	ClassifyError string    `gorm:"type:text" json:"classify_error"`        // Error message if classification failed
	ClassifiedAt  time.Time `json:"classified_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TrashRecord) TableName() string {
	return "trash_records"
}

// BeforeCreate hook to generate UUID if not set
func (t *TrashRecord) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
