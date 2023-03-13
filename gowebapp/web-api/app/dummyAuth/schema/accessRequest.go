package schema

import (
	// "time"

	// "log"
	// "web-api/app/dummyAuth/models"

	"github.com/google/uuid"
)

type AccessRequest struct {
	Base      Base      `gorm:"embedded"`
	ProfileId uuid.UUID `gorm:"type:uuid;not null;unique;indexed" json:"profileId"`
	EmailId   uuid.UUID `gorm:"type:uuid;not null;unique;indexed" json:"emailId"` //
	Approved  bool      `gorm:"default:false;not null" json:"approved,omitempty"` //
}
