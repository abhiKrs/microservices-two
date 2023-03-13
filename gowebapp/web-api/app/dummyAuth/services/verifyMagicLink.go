package services

import (
	"errors"
	"net/http"
	"time"

	constants "web-api/app/constants"
	dbConstants "web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (as *AuthService) VerifyMagicLink(req *models.SignInRequest) (*models.SignInResponse, error) {

	// var dbEmail schema.Email
	var err error

	if req.AuthType != uint8(constants.MagicLink) {
		log.ErrorLogger.Println("wrong function call")
		response := models.SignInResponse{IsSuccessful: false, Message: []string{http.StatusText(http.StatusInternalServerError)}}
		return &response, respond.ErrInternalServerError
	}

	tt, err := uuid.Parse(req.Credential)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		response := models.SignInResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	//  Get Magic Link
	dbMagicLink, err := as.dao.MagicLinkAccess.FindMagicLink(tt)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := models.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{"Token Not Available"}}
			return &response, err
		}
		response := models.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{err.Error()}}
		return &response, err
	}

	log.DebugLogger.Println(dbMagicLink)

	timeNow := time.Now()
	if dbMagicLink.ExpiryTime.Unix() < timeNow.Unix() {
		log.WarningLogger.Println("token expired")
		// http.Error(w, "token expired", http.StatusNotAcceptable)
		// Delete Magic Link
		as.dao.MagicLinkAccess.DeleteMagicLink(dbMagicLink)

		response := models.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{"token has expired!!!"}}
		return &response, respond.ErrTokenExpired
	}
	// dbProfile

	switch dbMagicLink.LinkType {
	case uint8(dbConstants.Register.EnumIndex()):
		// register magiclink route
		dbProfile, _, err := as.dao.VerifyRegisterLinkAndCreateProfile(dbMagicLink)

		if err != nil {
			log.ErrorLogger.Println(err.Error())
			response := models.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{err.Error()}}
			return &response, err
		}
		log.DebugLogger.Println(dbMagicLink.UserId.String())

		dbEmail, err := as.dao.EmailAccess.FindEmailbyId(dbMagicLink.EmailId)
		if err != nil {
			log.WarningLogger.Println(err)
		}
		dbTeam, err := as.dao.TeamAccess.GetByCreatorId(dbProfile.ID)
		if err != nil {
			log.WarningLogger.Println(err)
		}

		// dbAccessRequest, err := as.dao.AccessRequestAccess.FindAccessRequest(dbProfile.ID, dbEmail.ID, 1, 1)

		bodyPayload := models.ProfileBody{
			ProfileId:   dbProfile.Base.ID.String(),
			AccountId:   dbProfile.AccountId,
			Onboarded:   &dbProfile.Onboarded,
			ProfileBase: dbProfile.ProfileBase,
			Email:       dbEmail.Email,
			TeamId:      dbTeam.Base.ID.String(),
			// AccessApproved: dbAccessRequest.Approved,
		}

		response := models.SignInResponse{IsSuccessful: true, Code: int(dbMagicLink.LinkType), Email: dbEmail.Email, UserBody: bodyPayload, BearerToken: "bearer_" + dbProfile.Base.ID.String()}
		log.DebugLogger.Println(response)
		return &response, nil

	case uint8(dbConstants.Login), uint8(dbConstants.ResetPassword):
		//  login magiclink route
		// Get Profile
		dbProfile, err := as.dao.VerifyLoginLinkAndGetProfile(dbMagicLink)

		if err != nil {
			log.ErrorLogger.Println(err.Error())
			response := models.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{err.Error()}}
			return &response, err
		}

		dbEmail, err := as.dao.EmailAccess.FindEmailbyId(dbMagicLink.EmailId)
		if err != nil {
			log.WarningLogger.Println(err)
		}
		dbTeam, err := as.dao.TeamAccess.GetByCreatorId(dbProfile.ID)
		if err != nil {
			log.WarningLogger.Println(err)
		}
		dbAccessRequest, err := as.dao.AccessRequestAccess.FindAccessRequestByProfileId(dbProfile.ID)
		if err != nil {
			log.ErrorLogger.Println(err)
		}
		bodyPayload := models.ProfileBody{
			ProfileId:      dbProfile.Base.ID.String(),
			AccountId:      dbProfile.AccountId,
			Onboarded:      &dbProfile.Onboarded,
			ProfileBase:    dbProfile.ProfileBase,
			Email:          dbEmail.Email,
			AccessApproved: &dbAccessRequest.Approved,
			TeamId:         dbTeam.Base.ID.String(),
		}

		response := models.SignInResponse{IsSuccessful: true, Code: int(dbMagicLink.LinkType), Email: dbEmail.Email, UserBody: bodyPayload, BearerToken: "bearer_" + dbProfile.Base.ID.String()}
		log.DebugLogger.Println(response)
		return &response, nil

	default:
		response := models.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{http.StatusText(http.StatusInternalServerError)}}
		return &response, respond.ErrInternalServerError

	}

}
