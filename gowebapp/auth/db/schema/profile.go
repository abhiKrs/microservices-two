package schema

import (
	"github.com/google/uuid"
)

type Profile struct {
	Base                   //        Base        `gorm:"embedded"`
	UserId      uuid.UUID  `gorm:"type:uuid;not null;unique;indexed" json:"userId"`
	AccountId   *uuid.UUID `gorm:"type:uuid" json:"accountId,omitempty"`
	Onboarded   bool       `gorm:"default:false;not null"` // json:"onboarded,omitempty"
	ProfileBase            // ProfileBase `gorm:"embedded"`
}

type ProfileBase struct {
	FirstName    *string `json:"firstName,omitempty"`
	LastName     *string `json:"lastName,omitempty"`
	Title        *string `json:"title,omitempty"`
	Organization *string `json:"organization,omitempty"`
	Website      *string `json:"website,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	AvatarUrl    *string `json:"avatarUrl,omitempty"`
}

// func ProfileTable(env string) func(tx *gorm.DB) *gorm.DB {
// 	return func(tx *gorm.DB) *gorm.DB {
// 		log.Println(config.API().ENV)
// 		if env == constants.Test.String() {
// 			return tx.Table("test_profiles")
// 		} else if env == constants.Development.String() {
// 			return tx.Table("profiles")
// 		} else {
// 			log.Println(config.API().ENV)
// 			panic(errors.New("wrong environment"))
// 		}

// 		//   return tx.Table("profiles")
// 	}
// }
