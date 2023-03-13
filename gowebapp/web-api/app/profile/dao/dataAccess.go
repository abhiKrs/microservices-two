package dao

import (
	// "log"
	// "net/http"

	// "time"
	// "log"
	"web-api/app/config"
	// "web-api/app/dummyProfile/constants"
	"web-api/app/dummyAuth/dao/pgsql"
	teamPgsql "web-api/app/team/dao/pgsql"

	// "web-api/app/dummyAuth/models"
	// "web-api/app/dummyAuth/schema"

	// "web-api/app/dummyProfile/schema"
	// "web-api/app/dummyProfile/utility"

	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileDataAccess struct {
	cfg                 *config.Config
	db                  *gorm.DB
	EmailAccess         pgsql.EmailAccessOperator
	UserAccess          pgsql.UserAccessOperator
	CredentialAccess    pgsql.CredentialAccessOperator
	MagicLinkAccess     pgsql.MagicLinkAccessOperator
	AccessRequestAccess pgsql.AccessRequestAccessOperator
	ProfileAccess       pgsql.ProfileAccessOperator
	TeamAccess          teamPgsql.TeamAccessOperator
}

func New(db *gorm.DB, cfg *config.Config) *ProfileDataAccess {
	return &ProfileDataAccess{
		cfg:                 cfg,
		db:                  db,
		EmailAccess:         *pgsql.NewEmailAccess(db),
		UserAccess:          *pgsql.NewUserAccess(db),
		CredentialAccess:    *pgsql.NewcredentialAccess(db),
		AccessRequestAccess: *pgsql.NewAccessRequestAccess(db),
		MagicLinkAccess:     *pgsql.NewMagicLinkAccess(db),
		ProfileAccess:       *pgsql.NewProfileAccess(db),
		TeamAccess:          *teamPgsql.NewTeamAccess(db),
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
