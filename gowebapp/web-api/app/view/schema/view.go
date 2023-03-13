package schema

import (
	"time"

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

type View struct {
	Base          Base       `gorm:"embedded"`
	TeamId        uuid.UUID  `gorm:"type:uuid" json:"teamId"`
	StreamId      *uuid.UUID `gorm:"type:uuid" json:"streamId,omitempty"`
	SourcesFilter *[]string  `gorm:"type:text[]" json:"sourcesFilter,omitempty"`
	LevelFilter   *[]string  `gorm:"type:text[]" json:"levelFilter,omitempty"`
	StartDate     *time.Time `json:"startDate,omitempty"`
	EndDate       *time.Time `json:"endDate,omitempty"`
	Name          string     `gorm:"not null" json:"sourceType"`
}
