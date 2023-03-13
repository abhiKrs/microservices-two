package pgsql

import (
	"web-api/app/source/constants"
	"web-api/app/source/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SourceAccessOperator struct {
	db *gorm.DB
}

func NewSourceAccess(db *gorm.DB) *SourceAccessOperator {
	return &SourceAccessOperator{
		db: db,
	}
}

func (sao *SourceAccessOperator) Create(profileId uuid.UUID, name string, sType constants.SourceType, teamId uuid.UUID) (*schema.Source, error) {
	dbSource := schema.Source{ProfileId: profileId, Name: name, SourceType: sType, TeamId: teamId}

	result := sao.db.Create(&dbSource)

	return &dbSource, result.Error

}

func (sao *SourceAccessOperator) Update(dbModel *schema.Source, name string) (*schema.Source, error) {
	// var dbSource schema.Source
	// result := sao.db.First(&dbSource, id).Update(&schema.Source{SourceBase: *updateModel})
	result := sao.db.Model(&dbModel).Updates(schema.Source{Name: name})
	// result := sao.db.Model(&dbSource).Select("*").Updates(*updateModel)

	return dbModel, result.Error

}

func (sao *SourceAccessOperator) GetById(id uuid.UUID) (*schema.Source, error) {
	var dbSource schema.Source

	result := sao.db.First(&dbSource, id)

	return &dbSource, result.Error

}

func (sao *SourceAccessOperator) GetByProfileId(profileId uuid.UUID) (*[]schema.Source, error) {
	var dbSources []schema.Source

	result := sao.db.Where(&schema.Source{ProfileId: profileId}).Find(&dbSources)

	return &dbSources, result.Error

}

func (sao *SourceAccessOperator) DeleteById(id uuid.UUID) (*schema.Source, error) {
	var dbSource schema.Source

	result := sao.db.Delete(&dbSource, id)

	return &dbSource, result.Error

}
