package services

import (
	"web-api/app/constants"
	dbConstants "web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"

	// "web-api/app/dummyAuth/schema"
	"web-api/app/dummyAuth/utility"
	"web-api/app/utility/email"

	"github.com/google/uuid"
)

func (ps *ProfileService) ResetPassword(profileId uuid.UUID, password string) (*models.ProfileResponse, error) {

	log.DebugLogger.Println("started password reset service")

	dbProfile, err := ps.dao.ProfileAccess.GetById(profileId)

	// Get Email
	dbEmail, err := ps.dao.EmailAccess.FindPrimaryEmailbyUserId(dbProfile.UserId)
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	// Get Credential
	dbCredential, err := ps.dao.CredentialAccess.GetCredential(dbProfile.UserId, dbConstants.Custom, dbConstants.DBPassword)
	if err != nil {
		log.ErrorLogger.Println(err)
		log.InfoLogger.Println("password not reset")
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err

	} else {
		hashPass, err := utility.Hash(password)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
			return &response, err
		}

		// Update Old Credential
		dbCredential, err = ps.dao.CredentialAccess.UpdatePasswordCredential(dbCredential, string(hashPass), dbConstants.Bcrypt)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
			return &response, err
		}
	}

	successMessage := constants.PasswordReset
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
