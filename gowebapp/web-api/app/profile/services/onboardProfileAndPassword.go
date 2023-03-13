package services

import (
	"web-api/app/constants"
	"web-api/app/dummyAuth/models"
	"web-api/app/utility/email"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

func (ps *ProfileService) OnboardProfileAndPassword(profileId uuid.UUID, req *models.ProfileUpdateRequest) (*models.ProfileResponse, error) {

	var err error

	// Get Profile
	dbProfile, err := ps.dao.ProfileAccess.GetById(profileId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}
	// Check Not Onboarded
	if dbProfile.Onboarded {
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{"User has already been onboarded"}}
		return &response, respond.ErrBadRequest
	}

	// Define ProfileUpdate Model
	var updateModel = req.ProfileBase

	// Get Email(Primary)
	dbEmail, err := ps.dao.EmailAccess.FindPrimaryEmailbyUserId(dbProfile.UserId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	// Create hashPassword
	// hasPass, _ := utility.Hash(*req.Password)
	// strHasPass := string(hasPass)

	// Start Onboarding Transaction
	// dbProfile, dbAccessRequest, err := ps.dao.OnboardingTransaction(dbProfile, &updateModel, dbEmail.ID, &strHasPass)
	dbProfile, dbAccessRequest, err := ps.dao.OnboardingTransaction(dbProfile, &updateModel, dbEmail.ID)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	// utility.SendEmailviaAws(dbEmail.Email, )
	successMessage := constants.Onboarding
	// err = email.SendEmailviaDapr(dbEmail.Email, nil, constants.SuccessMessage, nil, &successMessage)
	err = email.SendEmail(dbEmail.Email, nil, constants.SuccessMessage, nil, &successMessage)
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	bodyPayload := models.ProfileBodyResponse{
		Profile:        *dbProfile,
		Email:          &dbEmail.Email,
		AccessApproved: dbAccessRequest.Approved,
	}
	response := models.ProfileResponse{IsSuccessful: true, UserBody: bodyPayload}
	return &response, nil

}
