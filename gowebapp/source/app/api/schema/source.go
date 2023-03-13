package schema

import (
	"github.com/google/uuid"
)

type Source struct {
	Base         Base      `gorm:"embedded"`
	ProfileId    uuid.UUID `gorm:"type:uuid;not null" json:"profileId"`
	Name         string    `gorm:"not null" json:"name"`
	PlatformType uint8     `gorm:"not null" json:"platformType"`
}
