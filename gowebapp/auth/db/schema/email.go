package schema

import (
	"time"

	"github.com/google/uuid"
)

type Email struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Email     string    `gorm:"size:100;not null;unique;indexed" json:"email"`
	Verified  bool      `gorm:"not null;" json:"verified"`
	UserId    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	Primary   bool      `gorm:"default:false;not null;" json:"primary"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// func EmailTable(env string) func(tx *gorm.DB) *gorm.DB {
// 	return func(tx *gorm.DB) *gorm.DB {
// 		log.Println(config.API().ENV)
// 		if env == constants.Test.String() {
// 			return tx.Table("test_emails")
// 		} else if env == constants.Development.String() {
// 			return tx.Table("emails")
// 		} else {
// 			log.Println(config.API().ENV)
// 			panic(errors.New("wrong environment"))
// 		}

// 		//   return tx.Table("emails")
// 	}
// }
