package services

import (
	// "errors"
	// "encoding/json"

	Constants "auth/constants"
	dbConstants "auth/db/constants"
	"auth/model"
	"auth/utility/respond"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (as *AuthService) VerifyMagicLink(req *model.SignInRequest) (*model.SignInResponse, error) {

	// var dbEmail schema.Email
	var err error

	if req.AuthType != uint8(Constants.MagicLink) {
		log.Println("wrong function call")
		response := model.SignInResponse{IsSuccessful: false, Message: []string{http.StatusText(http.StatusInternalServerError)}}
		return &response, respond.ErrInternalServerError
	}

	tt, err := uuid.Parse(req.Credential)
	if err != nil {
		log.Println(err.Error())
		response := model.SignInResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	//  Get Magic Link
	dbMagicLink, err := as.dao.MagicLinkAccess.FindMagicLink(tt)

	if err != nil {
		if errors.Is(err, respond.ErrNoRecord) {
			response := model.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{"Token Not Available"}}
			return &response, err
		}
		response := model.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{err.Error()}}
		return &response, err
	}

	log.Println(dbMagicLink)

	timeNow := time.Now()
	if dbMagicLink.ExpiryTime.Unix() < timeNow.Unix() {
		log.Println("token expired")
		// http.Error(w, "token expired", http.StatusNotAcceptable)
		// Delete Magic Link
		as.dao.MagicLinkAccess.DeleteMagicLink(dbMagicLink)

		response := model.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{"token has expired!!!"}}
		return &response, respond.ErrTokenExpired
	}
	// dbProfile

	switch dbMagicLink.LinkType {
	case uint8(dbConstants.Register.EnumIndex()):
		// register magiclink route
		dbProfile, _, err := as.dao.VerifyRegisterLinkAndCreateProfile(dbMagicLink)

		if err != nil {
			log.Println(err.Error())
			response := model.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{err.Error()}}
			return &response, err
		}
		log.Println(dbMagicLink.UserId.String())

		dbEmail, err := as.dao.EmailAccess.FindEmailbyId(dbMagicLink.EmailId)

		// dbAccessRequest, err := as.dao.AccessRequestAccess.FindAccessRequest(dbProfile.ID, dbEmail.ID, 1, 1)

		bodyPayload := model.ProfileBody{
			ProfileId:   dbProfile.Base.ID.String(),
			AccountId:   dbProfile.AccountId,
			Onboarded:   &dbProfile.Onboarded,
			ProfileBase: dbProfile.ProfileBase,
			Email:       dbEmail.Email,
			// AccessApproved: dbAccessRequest.Approved,
		}

		response := model.SignInResponse{IsSuccessful: true, Code: int(dbMagicLink.LinkType), Email: dbEmail.Email, UserBody: bodyPayload, BearerToken: "bearer_" + dbMagicLink.UserId.String()}
		log.Println(response)
		return &response, nil

	case uint8(dbConstants.Login), uint8(dbConstants.ResetPassword):
		//  login magiclink route
		// Get Profile
		dbProfile, err := as.dao.VerifyLoginLinkAndGetProfile(dbMagicLink)

		if err != nil {
			log.Println(err.Error())
			response := model.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{err.Error()}}
			return &response, err
		}

		dbEmail, err := as.dao.EmailAccess.FindEmailbyId(dbMagicLink.EmailId)
		dbAccessRequest, err := as.dao.AccessRequestAccess.FindAccessRequestByProfileId(dbProfile.ID)
		if err != nil {
			log.Println(err)
		}
		bodyPayload := model.ProfileBody{
			ProfileId:      dbProfile.Base.ID.String(),
			AccountId:      dbProfile.AccountId,
			Onboarded:      &dbProfile.Onboarded,
			ProfileBase:    dbProfile.ProfileBase,
			Email:          dbEmail.Email,
			AccessApproved: &dbAccessRequest.Approved,
		}

		response := model.SignInResponse{IsSuccessful: true, Code: int(dbMagicLink.LinkType), Email: dbEmail.Email, UserBody: bodyPayload, BearerToken: "bearer_" + dbMagicLink.UserId.String()}
		log.Println(response)
		return &response, nil

	default:
		response := model.SignInResponse{IsSuccessful: false, Code: int(dbMagicLink.LinkType), Message: []string{http.StatusText(http.StatusInternalServerError)}}
		return &response, respond.ErrInternalServerError

	}

}
