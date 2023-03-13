package services

import (
	"errors"
	"net/http"

	"web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/models"
	"web-api/app/dummyAuth/utility"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"

	"gorm.io/gorm"
)

func (as *AuthService) VerifySigninPasswordAndGetProfile(req *models.SignInRequest) (*models.SignInResponse, error) {
	log.DebugLogger.Println("inside pass and get profile")
	var err error

	// get email
	dbEmail, err := as.dao.EmailAccess.FindEmailbyStrAdd(*req.Email)
	if err != nil {
		log.ErrorLogger.Println("email not found")
		response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
		return &response, err
	}

	// get credential
	dbCredential, err := as.dao.CredentialAccess.GetCredential(dbEmail.UserId, constants.Custom, constants.DBPassword)

	if err != nil {
		log.ErrorLogger.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
			return &response, err
		}
		response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
		return &response, err
	}

	passAlgo := *dbCredential.PasswordAlgo
	// log.DebugLogger.Println(passAlgo)

	switch passAlgo {
	case "bcrypt":
		err = utility.VerifyPassword(*dbCredential.PasswordHash, req.Credential)

		if err != nil {
			log.ErrorLogger.Println("wrong password")
			response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{"Password didn't match!!!"}}
			return &response, err
		}
	default:
		response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{http.StatusText(http.StatusInternalServerError)}}
		return &response, respond.ErrInternalServerError
	}

	// get profile
	dbProfile, err := as.dao.ProfileAccess.GetByUserId(dbCredential.UserId)
	if err != nil {
		log.ErrorLogger.Println("credentials not found")
		response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
		return &response, err
	}

	// get Access Request
	dbAccessRequests, err := as.dao.AccessRequestAccess.FindAccessRequestByProfileId(dbProfile.ID)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	dbTeam, err := as.dao.TeamAccess.GetByCreatorId(dbProfile.ID)
	if err != nil {
		log.WarningLogger.Println(err)
	}

	bodyPayload := models.ProfileBody{
		ProfileBase:    dbProfile.ProfileBase,
		ProfileId:      dbProfile.Base.ID.String(),
		AccountId:      dbProfile.AccountId,
		Onboarded:      &dbProfile.Onboarded,
		AccessApproved: &dbAccessRequests.Approved,
		Email:          dbEmail.Email,
		TeamId:         dbTeam.Base.ID.String(),
	}

	//  TODO

	response := models.SignInResponse{IsSuccessful: true, Code: 2, UserBody: bodyPayload, BearerToken: "bearer_" + dbProfile.Base.ID.String()}
	return &response, nil

}
