package dao

import (
	"web-api/app/config"
	authPgsql "web-api/app/dummyAuth/dao/pgsql"
	sourcePgsql "web-api/app/source/dao/pgsql"
	teamPgsql "web-api/app/team/dao/pgsql"
	"web-api/app/view/dao/pgsql"

	"gorm.io/gorm"
)

type ViewDataAccess struct {
	cfg                 *config.Config
	db                  *gorm.DB
	ViewAccess          pgsql.ViewAccessOperator
	SourceAccess        sourcePgsql.SourceAccessOperator
	EmailAccess         authPgsql.EmailAccessOperator
	UserAccess          authPgsql.UserAccessOperator
	CredentialAccess    authPgsql.CredentialAccessOperator
	MagicLinkAccess     authPgsql.MagicLinkAccessOperator
	AccessRequestAccess authPgsql.AccessRequestAccessOperator
	ProfileAccess       authPgsql.ProfileAccessOperator
	TeamAccess          teamPgsql.TeamAccessOperator
}

func New(db *gorm.DB, cfg *config.Config) *ViewDataAccess {
	return &ViewDataAccess{
		cfg:                 cfg,
		db:                  db,
		ViewAccess:          *pgsql.NewViewAccess(db),
		SourceAccess:        *sourcePgsql.NewSourceAccess(db),
		EmailAccess:         *authPgsql.NewEmailAccess(db),
		UserAccess:          *authPgsql.NewUserAccess(db),
		CredentialAccess:    *authPgsql.NewcredentialAccess(db),
		AccessRequestAccess: *authPgsql.NewAccessRequestAccess(db),
		MagicLinkAccess:     *authPgsql.NewMagicLinkAccess(db),
		ProfileAccess:       *authPgsql.NewProfileAccess(db),
		TeamAccess:          *teamPgsql.NewTeamAccess(db),
	}
}
