package dao

import (
	// "log"
	// "net/http"

	// "time"
	// "log"
	"web-api/app/config"
	// "web-api/app/dummyProfile/constants"
	authPgsql "web-api/app/dummyAuth/dao/pgsql"
	sourcePgsql "web-api/app/source/dao/pgsql"
	"web-api/app/team/dao/pgsql"

	// "web-api/app/dummyAuth/models"
	// "web-api/app/dummyAuth/schema"

	// "web-api/app/dummyProfile/schema"
	// "web-api/app/dummyProfile/utility"

	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamDataAccess struct {
	cfg                 *config.Config
	db                  *gorm.DB
	EmailAccess         authPgsql.EmailAccessOperator
	UserAccess          authPgsql.UserAccessOperator
	CredentialAccess    authPgsql.CredentialAccessOperator
	MagicLinkAccess     authPgsql.MagicLinkAccessOperator
	AccessRequestAccess authPgsql.AccessRequestAccessOperator
	ProfileAccess       authPgsql.ProfileAccessOperator
	SourceAccess        sourcePgsql.SourceAccessOperator
	TeamAccess          pgsql.TeamAccessOperator
}

func New(db *gorm.DB, cfg *config.Config) *TeamDataAccess {
	return &TeamDataAccess{
		cfg:                 cfg,
		db:                  db,
		SourceAccess:        *sourcePgsql.NewSourceAccess(db),
		EmailAccess:         *authPgsql.NewEmailAccess(db),
		UserAccess:          *authPgsql.NewUserAccess(db),
		CredentialAccess:    *authPgsql.NewcredentialAccess(db),
		AccessRequestAccess: *authPgsql.NewAccessRequestAccess(db),
		MagicLinkAccess:     *authPgsql.NewMagicLinkAccess(db),
		ProfileAccess:       *authPgsql.NewProfileAccess(db),
		TeamAccess:          *pgsql.NewTeamAccess(db),
	}
}

// func (pda *ProfileDataAccess) UpdateProfile(updateRequestModel *models.ProfileUpdateRequest) (*schema.Profile, error) {
// 	// var dbUser *schema.User
// 	// var dbEmail *schema.Email
// 	var dbCredential *schema.Credential
// 	var password = updateRequestModel.Password

// 	// first get credential

// 	dbCredential = pda.CredentialAccess.GetCredentials()

// 	err = pda.db.Transaction(func(tx *gorm.DB) error {
// 		var err error
// 		dbCredential = &schema.Credential{

// 		}

// 	})
// 	// 	dbUser = &schema.User{
// 	// 		Base: schema.Base{},
// 	// 		Role: uint8(constants.Users),
// 	// 	}

// 	// 	result := tx.Create(&dbUser)
// 	// 	if result.Error != nil {
// 	// 		return err
// 	// 	}

// 	// 	dbEmail = &schema.Email{Email: emailNew, UserId: dbUser.Base.ID, Primary: true}
// 	// 	result = tx.Create(&dbEmail)
// 	// 	if result.Error != nil {
// 	// 		return err
// 	// 	}

// 	// 	dbMagicLink, err = au.MagicLinkAccess.CreateMagicLink(constants.Register, dbEmail.ID, dbUser.Base.ID, tx)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}

// 	// 	tokenStr := dbMagicLink.Token.String()
// 	// 	recepient := emailNew

// 	// 	res := utility.SendEmailviaAws(recepient, tokenStr, constants.Register, au.cfg)

// 	// 	if res != nil {
// 	// 		log.Println(res.Error())
// 	// 		return err
// 	// 	}
// 	// 	return nil
// 	// })
// 	// if err != nil {
// 	// 	log.Println("errorn in transaction")
// 	// 	return nil, nil, err
// 	// }

// }
