package schema

import (
	"errors"

	"web-api/app/config"
	"web-api/app/constants"
	log "web-api/app/utility/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
	// "gorm.io/gorm"
)

type Credential struct {
	Base         Base      `gorm:"embedded"`
	UserId       uuid.UUID `gorm:"type:uuid;not null"`
	AuthType     uint8     `gorm:"not null" json:"authType"`
	ProviderType uint8     `gorm:"not null" json:"providerType"`
	Identifier   string    `gorm:"size:100;not null" json:"identifier"`
	PasswordHash *string   `json:"passwordHash"`
	PasswordAlgo *string   `json:"passwordAlgo"`
}

func CredentialTable(env string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		log.DebugLogger.Println(config.API().ENV)
		if env == constants.Test.String() {
			return tx.Table("test_credentials")
		} else if env == constants.Development.String() {
			return tx.Table("credentials")
		} else {
			log.DebugLogger.Println(config.API().ENV)
			panic(errors.New("wrong environment"))
		}

		//   return tx.Table("credentials")
	}
}
