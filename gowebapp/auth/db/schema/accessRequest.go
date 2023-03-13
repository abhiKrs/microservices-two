package schema

import (
	// "time"

	// "log"
	// "web-api/app/dummyAuth/models"

	"github.com/google/uuid"
)

type AccessRequest struct {
	Base      Base
	ProfileId uuid.UUID `json:"profileId"`
	EmailId   uuid.UUID `json:"emailId"`            //
	Approved  bool      `json:"approved,omitempty"` //
}
