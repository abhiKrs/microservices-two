package schema

import (
	"time"

	// "github.com/badoux/checkmail"
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
