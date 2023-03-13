package dao

import (
	"web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/schema"
	log "web-api/app/utility/logger"

	"gorm.io/gorm"
)

func (au *AuthDataAccess) CreateNewUserEmailAndGoogleCredential(emailNew string, verified bool, primary bool, identifier string, authType constants.DBAuthType, providerType constants.DBProviderType) (*schema.Email, *schema.User, error) {
	var dbUser *schema.User
	var dbEmail *schema.Email

	err := au.db.Transaction(func(tx *gorm.DB) error {
		var err error
		dbEmail, dbUser, err = au.CreateNewUser(emailNew, verified, primary)
		if err != nil {
			log.ErrorLogger.Println(err)
			return err
		}

		// Create Credential
		dbCredential, err := au.CredentialAccess.CreateGoogleCredential(dbUser.Base.ID, identifier, authType, providerType)
		if err != nil {
			log.ErrorLogger.Println(err)
			return err
		}
		log.DebugLogger.Println(dbCredential)

		return nil
	})
	if err != nil {
		log.ErrorLogger.Println(err)
		log.ErrorLogger.Println("errorn in transaction")
		return nil, nil, err
	}
	return dbEmail, dbUser, nil
}
