package pgsql

import (
	"web-api/app/dummyAuth/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileAccessOperator struct {
	db *gorm.DB
}

func NewProfileAccess(db *gorm.DB) *ProfileAccessOperator {
	return &ProfileAccessOperator{
		db: db,
	}
}

func (uao *ProfileAccessOperator) Create(userId uuid.UUID) (*schema.Profile, error) {
	dbProfile := schema.Profile{UserId: userId, Onboarded: false}

	result := uao.db.Create(&dbProfile)

	return &dbProfile, result.Error

}

func (uao *ProfileAccessOperator) Update(dbModel *schema.Profile, updateModel *schema.ProfileBase) (*schema.Profile, error) {
	// var dbProfile schema.Profile
	// result := uao.db.First(&dbProfile, id).Update(&schema.Profile{ProfileBase: *updateModel})
	result := uao.db.Model(&dbModel).Updates(schema.Profile{ProfileBase: *updateModel, Onboarded: true})
	// result := uao.db.Model(&dbProfile).Select("*").Updates(*updateModel)

	return dbModel, result.Error

}

func (uao *ProfileAccessOperator) GetById(id uuid.UUID) (*schema.Profile, error) {
	var dbProfile schema.Profile

	result := uao.db.First(&dbProfile, id)

	return &dbProfile, result.Error

}

func (uao *ProfileAccessOperator) GetByUserId(userId uuid.UUID) (*schema.Profile, error) {
	var dbProfile schema.Profile

	result := uao.db.Where(&schema.Profile{UserId: userId}).First(&dbProfile)

	return &dbProfile, result.Error

}
