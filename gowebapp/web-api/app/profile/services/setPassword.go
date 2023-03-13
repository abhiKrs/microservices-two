package services

import (
	"errors"

	"web-api/app/constants"
	dbConstants "web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"

	// "web-api/app/dummyAuth/schema"
	"web-api/app/dummyAuth/utility"
	"web-api/app/utility/email"
	"web-api/app/utility/respond"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ps *ProfileService) SetPassword(profileId uuid.UUID, password string) (*models.ProfileResponse, error) {

	log.InfoLogger.Println("started password reset service")

	dbProfile, err := ps.dao.ProfileAccess.GetById(profileId)

	// Get Email
	dbEmail, err := ps.dao.EmailAccess.FindPrimaryEmailbyUserId(dbProfile.UserId)
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	// Get Credential
	_, err = ps.dao.CredentialAccess.GetCredential(dbProfile.UserId, dbConstants.Custom, dbConstants.DBPassword)
	// log.Println()
	if err != nil {
		log.ErrorLogger.Println(err)

		// Create New Credential(password)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hasPass, _ := utility.Hash(password)
			strHasPass := string(hasPass)
			_, err = ps.dao.CredentialAccess.CreatePasswordCredential(
				dbProfile.UserId,
				&strHasPass,
				dbConstants.Bcrypt,
				dbEmail.ID.String(),
				dbConstants.Custom,
				dbConstants.DBPassword,
			)
		} else {
			log.InfoLogger.Println("password not created")
			response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
			return &response, err
		}
	} else {
		log.InfoLogger.Println("password credential already exists")
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{"method not allowed. Password set already"}}
		return &response, respond.ErrInvalidRequest
	}

	successMessage := constants.PasswordSet
	// err = email.SendEmailviaDapr(dbEmail.Email, nil, constants.SuccessMessage, nil, &successMessage)
	err = email.SendEmail(dbEmail.Email, nil, constants.SuccessMessage, nil, &successMessage)

	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	log.DebugLogger.Println("finished")
	bodyPayload := models.ProfileBodyResponse{
		Profile: *dbProfile,
		Email:   &dbEmail.Email,
	}
	response := models.ProfileResponse{IsSuccessful: true, UserBody: bodyPayload}
	return &response, nil
}
