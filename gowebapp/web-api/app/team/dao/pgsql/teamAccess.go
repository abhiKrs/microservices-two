package pgsql

import (
	"web-api/app/team/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamAccessOperator struct {
	db *gorm.DB
}

func NewTeamAccess(db *gorm.DB) *TeamAccessOperator {
	return &TeamAccessOperator{
		db: db,
	}
}

func (sao *TeamAccessOperator) Create(profileId uuid.UUID, name string) (*schema.Team, error) {
	dbTeam := schema.Team{CreatedBy: profileId, Name: name}

	result := sao.db.Create(&dbTeam)

	return &dbTeam, result.Error

}

// func (sao *TeamAccessOperator) Update(dbModel *schema.Team, updateModel *schema.TeamBase) (*schema.Team, error) {
// 	// var dbTeam schema.Team
// 	// result := sao.db.First(&dbTeam, id).Update(&schema.Team{TeamBase: *updateModel})
// 	result := sao.db.Model(&dbModel).Updates(schema.Team{TeamBase: *updateModel, Onboarded: true})
// 	// result := sao.db.Model(&dbTeam).Select("*").Updates(*updateModel)

// 	return dbModel, result.Error

// }

func (sao *TeamAccessOperator) GetById(id uuid.UUID) (*schema.Team, error) {
	var dbTeam schema.Team

	result := sao.db.First(&dbTeam, id)

	return &dbTeam, result.Error

}

func (sao *TeamAccessOperator) GetByCreatorId(profileId uuid.UUID) (*schema.Team, error) {
	var dbTeams schema.Team

	result := sao.db.Where(&schema.Team{CreatedBy: profileId}).First(&dbTeams)

	return &dbTeams, result.Error

}
