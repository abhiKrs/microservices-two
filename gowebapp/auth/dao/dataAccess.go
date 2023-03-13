package dao

import (
	"log"
	"time"

	// "net/http"

	// "time"
	"auth/config"
	constants "auth/constants"
	dbConstants "auth/constants"
	"auth/dao/pgsql"
	"auth/db/schema"
	"auth/utility/email"
	"auth/utility/respond"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthDataAccess struct {
	cfg                 *config.Config
	db                  *gorm.DB
	EmailAccess         pgsql.EmailAccessOperator
	UserAccess          pgsql.UserAccessOperator
	CredentialAccess    pgsql.CredentialAccessOperator
	MagicLinkAccess     pgsql.MagicLinkAccessOperator
	ProfileAccess       pgsql.ProfileAccessOperator
	AccessRequestAccess pgsql.AccessRequestAccessOperator
}

func New(db *gorm.DB, cfg *config.Config) *AuthDataAccess {
	return &AuthDataAccess{
		cfg:                 cfg,
		db:                  db,
		EmailAccess:         *pgsql.NewEmailAccess(db),
		UserAccess:          *pgsql.NewUserAccess(db),
		CredentialAccess:    *pgsql.NewcredentialAccess(db),
		MagicLinkAccess:     *pgsql.NewMagicLinkAccess(db),
		AccessRequestAccess: *pgsql.NewAccessRequestAccess(db),
		ProfileAccess:       *pgsql.NewProfileAccess(db),
		// EmailAccess:      *pgsql.NewEmailAccess(db.Scopes(schema.EmailTable(cfg.Api.ENV))),
		// UserAccess:       *pgsql.NewUserAccess(db.Scopes(schema.UserTable(cfg.Api.ENV))),
		// CredentialAccess: *pgsql.NewcredentialAccess(db.Scopes(schema.CredentialTable(cfg.Api.ENV))),
		// MagicLinkAccess:  *pgsql.NewMagicLinkAccess(db.Scopes(schema.MagicLinkTable(cfg.Api.ENV))),
		// ProfileAccess:    *pgsql.NewProfileAccess(db.Scopes(schema.ProfileTable(cfg.Api.ENV))),
	}
}

func (au *AuthDataAccess) CreateNewUser(emailNew string, verified bool, primary bool) (*schema.Email, *schema.User, error) {
	var dbUser *schema.User
	var dbEmail *schema.Email
	err := au.db.Transaction(func(tx *gorm.DB) error {
		var err error
		dbUser = &schema.User{
			Base: schema.Base{},
			Role: uint8(dbConstants.Users),
		}

		result := tx.Create(&dbUser)
		// result := tx.Scopes(schema.UserTable(au.cfg.Api.ENV)).Create(&dbUser)
		if result.Error != nil {
			return err
		}

		dbEmail = &schema.Email{Email: emailNew, UserId: dbUser.Base.ID, Primary: primary, Verified: verified}
		result = tx.Create(&dbEmail)
		// result = tx.Scopes(schema.EmailTable(au.cfg.Api.ENV)).Create(&dbEmail)
		if result.Error != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("errorn in transaction")
		return nil, nil, err
	}
	return dbEmail, dbUser, nil
}

func (au *AuthDataAccess) CreateNewUserAndSendMagicLink(emailNew string) (*schema.Email, *schema.MagicLink, error) {
	var dbUser *schema.User
	var dbEmail *schema.Email
	var dbMagicLink *schema.MagicLink
	err := au.db.Transaction(func(tx *gorm.DB) error {
		var err error
		dbUser = &schema.User{
			Base: schema.Base{},
			Role: uint8(dbConstants.Users),
		}

		result := tx.Create(&dbUser)
		// result := tx.Scopes(schema.UserTable(au.cfg.Api.ENV)).Create(&dbUser)
		if result.Error != nil {
			return err
		}

		dbEmail = &schema.Email{Email: emailNew, UserId: dbUser.Base.ID, Primary: true}
		result = tx.Create(&dbEmail)
		// result = tx.Scopes(schema.EmailTable(au.cfg.Api.ENV)).Create(&dbEmail)
		if result.Error != nil {
			return err
		}

		token := uuid.New()

		// var dbMagicLink schema.MagicLink

		dbMagicLink = &schema.MagicLink{
			Token:      token,
			UserId:     dbUser.Base.ID,
			EmailId:    dbEmail.ID,
			LinkType:   uint8(dbConstants.Register),
			CreatedAt:  time.Now(),
			ExpiryTime: time.Now().Add(1 * time.Hour),
		}

		// result = tx.Scopes(schema.UserTable(au.cfg.Api.ENV)).Create(&dbMagicLink)
		result = tx.Create(&dbMagicLink)
		if result.Error != nil {
			return result.Error
		}

		tokenStr := dbMagicLink.Token.String()
		recepient := emailNew
		linkType := dbConstants.Register

		res := email.SendEmailviaAws(recepient, &tokenStr, constants.MagicLink, &linkType, nil, au.cfg)

		if res != nil {
			log.Println(res.Error())
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("errorn in transaction")
		return nil, nil, err
	}
	return dbEmail, dbMagicLink, nil
}

func (au *AuthDataAccess) ResendAuthLink(dbEmail *schema.Email, linkType dbConstants.DBMagicLinkType) (*schema.MagicLink, error) {
	var dbUser *schema.User
	var dbMagicLink *schema.MagicLink
	var errr error
	err := au.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Find(&dbUser, dbEmail.UserId)

		if result.Error != nil {
			return result.Error
		}

		token := uuid.New()

		dbMagicLink = &schema.MagicLink{
			Token:      token,
			UserId:     dbUser.Base.ID,
			EmailId:    dbEmail.ID,
			LinkType:   uint8(linkType),
			CreatedAt:  time.Now(),
			ExpiryTime: time.Now().Add(1 * time.Hour),
		}

		result = tx.Create(&dbMagicLink)
		if result.Error != nil {
			return result.Error
		}

		tokenStr := dbMagicLink.Token.String()
		recipient := dbEmail.Email

		errr = email.SendEmailviaAws(recipient, &tokenStr, constants.MagicLink, &linkType, nil, au.cfg)

		if errr != nil {
			log.Println(errr)
			return errr
		}
		return nil

	})

	if err != nil {
		return dbMagicLink, err
	}
	return dbMagicLink, nil
}

func (au *AuthDataAccess) VerifyRegisterLinkAndCreateProfile(dbMagicLink *schema.MagicLink) (*schema.Profile, string, error) {
	// var email *string
	var dbProfile *schema.Profile
	if dbMagicLink.LinkType != uint8(dbConstants.Register) {
		log.Println("wrong function call")
		return dbProfile, "", respond.ErrInternalServerError
	}
	var dbUser *schema.User
	var dbEmail *schema.Email

	err := au.db.Transaction(func(tx *gorm.DB) error {
		//  get user
		result := tx.First(&dbUser, dbMagicLink.UserId)
		if result.Error != nil {
			return result.Error
		}
		// update email
		result = tx.First(&dbEmail, dbMagicLink.EmailId).Updates(schema.Email{Verified: true, Primary: true})
		if result.Error != nil {
			return result.Error
		}
		// create profile
		dbProfile = &schema.Profile{UserId: dbUser.Base.ID, Onboarded: false}
		result = tx.Create(&dbProfile)
		if result.Error != nil {
			return result.Error
		}
		// delete magic link
		result = tx.Delete(&dbMagicLink)
		if result.Error != nil {
			return result.Error
		}
		return nil

	})

	if err != nil {
		log.Println(err.Error())
		return dbProfile, dbEmail.Email, err
	}

	return dbProfile, dbEmail.Email, nil

}

func (au *AuthDataAccess) VerifyLoginLinkAndGetProfile(dbMagicLink *schema.MagicLink) (*schema.Profile, error) {
	var dbProfile *schema.Profile

	// if dbMagicLink.LinkType != uint8(dbconstants.Login) {
	// 	log.Println("wrong function call")
	// 	return dbProfile, respond.ErrInternalServerError
	// }
	var dbEmail *schema.Email

	// get email
	result := au.db.First(&dbEmail, dbMagicLink.EmailId)
	if result.Error != nil {
		return dbProfile, result.Error
	}

	if !dbEmail.Verified {
		log.Println("called verifyLogin when email not verified!!!!")
		return dbProfile, respond.ErrBadRequest
	}

	// delete magic link
	result = au.db.Delete(&dbMagicLink)
	if result.Error != nil {
		return dbProfile, result.Error
	}

	// get profile
	dbProfile, err := au.ProfileAccess.GetByUserId(dbMagicLink.UserId)
	if err != nil {
		return dbProfile, err
	}
	// result = au.db.Delete(&dbMagicLink)
	if err != nil {
		return dbProfile, err
		// *dbProfile = schema.Profile{UserId: dbMagicLink.UserId, ProfileBase: schema.ProfileBase{Onboarded: false}}
		// result = au.db.Create(&dbProfile)
		// if result.Error != nil {
		// 	return dbProfile, result.Error
		// }
		// return dbProfile, nil
	}
	return dbProfile, nil

}

func (au *AuthDataAccess) VerifyPasswordLink(dbMagicLink *schema.MagicLink) (*schema.MagicLink, error) {

	if dbMagicLink.LinkType != uint8(dbConstants.ResetPassword) {
		log.Println("wrong function call")
		return dbMagicLink, respond.ErrInternalServerError
	}
	var dbEmail *schema.Email

	// get email
	result := au.db.First(&dbEmail, dbMagicLink.EmailId)
	// tx.First()
	if result.Error != nil {
		return dbMagicLink, result.Error
	}

	if !dbEmail.Verified {
		log.Println("called verifyLogin when email not verified!!!!")
		return dbMagicLink, respond.ErrBadRequest
	}
	// delete magic link
	result = au.db.Delete(&dbMagicLink)
	if result.Error != nil {
		return dbMagicLink, result.Error
	}
	return dbMagicLink, nil

}
