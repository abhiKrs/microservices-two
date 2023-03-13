package dao

import (
	"web-api/app/config"
	"web-api/app/source/dao/pgsql"

	"gorm.io/gorm"
)

type SourceDataAccess struct {
	cfg          *config.Config
	SourceAccess pgsql.SourceAccessOperator
}

func New(db *gorm.DB, cfg *config.Config) *SourceDataAccess {
	return &SourceDataAccess{
		cfg:          cfg,
		SourceAccess: *pgsql.NewSourceAccess(db),
	}
}
