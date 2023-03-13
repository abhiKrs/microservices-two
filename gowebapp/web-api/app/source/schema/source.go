package schema

import (
	"time"
	"web-api/app/source/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	// gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key;autoIncrement:false" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"` // json:"createdAt"
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"` // json:"updatedAt"
	DeletedAt gorm.DeletedAt
}

type Source struct {
	Base
	ProfileId  uuid.UUID            `gorm:"type:uuid;not null" json:"profileId"`
	TeamId     uuid.UUID            `gorm:"type:uuid;not null" json:"teamId"`
	Name       string               `gorm:"not null" json:"name"`
	SourceType constants.SourceType `gorm:"not null" json:"sourceType"`
}
