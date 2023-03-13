package schema

import (
	"encoding/json"
	"errors"
	"io"

	"web-api/app/config"
	"web-api/app/constants"
	log "web-api/app/utility/logger"

	"gorm.io/gorm"
)

type User struct {
	Base Base  `gorm:"embedded"`
	Role uint8 `gorm:"size:255" json:"role"`
}

type CreateRequest struct {
	Role      string `json:"role" validate:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateRequest struct {
	ID   int    `json:"id"`
	Role string `json:"role,omitempty"`
	// Onboarded
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

func (r *UpdateRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

func (r *CreateRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

func UserTable(env string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		log.DebugLogger.Println(config.API().ENV)
		if env == constants.Test.String() {
			return tx.Table("test_users")
		} else if env == constants.Development.String() {
			return tx.Table("users")
		} else {
			log.DebugLogger.Println(config.API().ENV)
			panic(errors.New("wrong environment"))
		}

		//   return tx.Table("users")
	}
}
