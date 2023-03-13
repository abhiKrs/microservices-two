package schema

import (
	"time"

	"github.com/google/uuid"
)

type MagicLink struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Token      uuid.UUID `gorm:"type:uuid;not null;unique" json:"token"`
	UserId     uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	EmailId    uuid.UUID `gorm:"type:uuid;not null" json:"emailId"`
	LinkType   uint8     `gorm:"not null" json:"linkType"`
	ExpiryTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"expiryTime"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// func MagicLinkTable(env string) func(tx *gorm.DB) *gorm.DB {
// 	return func(tx *gorm.DB) *gorm.DB {
// 		log.Println(config.API().ENV)
// 		if env == constants.Test.String() {
// 			return tx.Table("test_magic_links")
// 		} else if env == constants.Development.String() {
// 			return tx.Table("magic_links")
// 		} else {
// 			log.Println(config.API().ENV)
// 			panic(errors.New("wrong environment"))
// 		}

// 		//   return tx.Table("magic_links")
// 	}
// }
