package pgsql

import (
	"web-api/app/view/models"
	"web-api/app/view/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ViewAccessOperator struct {
	db *gorm.DB
}

func NewViewAccess(db *gorm.DB) *ViewAccessOperator {
	return &ViewAccessOperator{
		db: db,
	}
}

func (sao *ViewAccessOperator) Create(teamId uuid.UUID, name string, levels *[]string, sources *[]string, date *models.DateInterval, streamId *uuid.UUID) (*schema.View, error) {
	dbView := schema.View{
		TeamId: teamId,
		Name:   name,
	}
	if date != nil {
		if date.StartDate != nil {
			dbView.StartDate = date.StartDate
		}
		if date.EndDate != nil {
			dbView.EndDate = date.EndDate
		}
	}
	if streamId != nil {
		dbView.StreamId = streamId
	}
	if sources != nil {
		dbView.SourcesFilter = sources
	}
	if levels != nil {
		dbView.LevelFilter = levels
	}

	result := sao.db.Create(&dbView)

	return &dbView, result.Error

}

func (sao *ViewAccessOperator) Update(dbModel *schema.View, pModel models.UpdateViewRequest) (*schema.View, error) {
	// var dbView schema.View
	updateModel := schema.View{
		Name:          *pModel.Name,
		SourcesFilter: pModel.SourcesFilter,
		LevelFilter:   pModel.LevelFilter,
		StartDate:     pModel.DateFilter.StartDate,
		EndDate:       pModel.DateFilter.EndDate,
	}
	if pModel.StreamId != nil {
		sid, _ := uuid.Parse(*pModel.StreamId)
		updateModel.StreamId = &sid
	}

	// result := sao.db.First(&dbView, id).Update(&schema.View{ViewBase: *updateModel})
	result := sao.db.Model(&dbModel).Updates(updateModel)
	// result := sao.db.Model(&dbView).Select("*").Updates(*updateModel)

	return dbModel, result.Error

}

func (sao *ViewAccessOperator) GetById(id uuid.UUID) (*schema.View, error) {
	var dbView schema.View

	result := sao.db.First(&dbView, id)

	return &dbView, result.Error

}

func (sao *ViewAccessOperator) GetByTeamId(teamId uuid.UUID) (*[]schema.View, error) {
	var dbViews []schema.View

	result := sao.db.Where(&schema.View{TeamId: teamId}).Find(&dbViews)

	return &dbViews, result.Error

}

func (sao *ViewAccessOperator) DeleteById(id uuid.UUID) (*schema.View, error) {
	var dbView schema.View
	result := sao.db.Delete(&dbView, id)
	return &dbView, result.Error
}
