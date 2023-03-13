package services

import (
	"context"
	"errors"
	"fmt"

	"web-api/app/config"
	"web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/models"
	"web-api/app/dummyAuth/schema"
	"web-api/app/utility/kafka"
	log "web-api/app/utility/logger"

	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

func (as *AuthService) VerifyGoogleSigninAndGetProfile(req *models.SignInRequest) (*models.SignInResponse, error) {
	log.InfoLogger.Println("inside google sign in and get profile")
	var err error
	log.DebugLogger.Println(req)
	var dbEmail *schema.Email
	var dbUser *schema.User
	var dbProfile *schema.Profile
	var dbAccessRequest *schema.AccessRequest
	// var dbCredential
	// func Validate(ctx context.Context, idToken string, audience string) (*Payload, error)
	ctx := context.Background()

	payload, err := idtoken.Validate(ctx, req.Credential, config.GOOGLE().CLIENT_ID)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{"illegal idToken"}}
		return &response, err
	}

	curEmail := fmt.Sprint(payload.Claims["email"])
	identifier := fmt.Sprint(payload.Claims["sub"])

	// Get Credentials
	// var dbCredential schema.Credential
	authType := constants.Custom
	providerType := constants.DBGoogle

	dbCredential, err := as.dao.CredentialAccess.GetCredentialbyIdentifier(identifier, authType, providerType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.ErrorLogger.Println(err)
			// create verified email and create a new account
			// get dbEmail
			verified := true
			dbEmail, err = as.dao.EmailAccess.FindEmailbyStrAdd(curEmail)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					log.ErrorLogger.Println(err)
					// no email entry
					// create new user, verified and primary email entry, create credential

					primary := true
					dbEmail, dbUser, err = as.dao.CreateNewUserEmailAndGoogleCredential(curEmail, verified, primary, identifier, authType, providerType)
					if err != nil {
						log.ErrorLogger.Println(err)
						response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{"illegal idToken"}}
						return &response, err
					}

					dbProfile, err = as.dao.ProfileAccess.Create(dbUser.Base.ID)
					if err != nil {
						log.ErrorLogger.Println(err)
						response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
						return &response, err
					}

					dbTeam, err := as.dao.TeamAccess.Create(dbProfile.ID, "My Team")
					if err != nil {
						log.ErrorLogger.Println(err)
						// return result.Error
					}
					// kafka-topic for team
					if dbTeam.Base.ID != uuid.Nil {
						teamTopic := "team_topic_" + dbTeam.Base.ID.String()
						go kafka.CreateTopic(teamTopic)
					}
					// temporary=================================================
					if dbProfile.ID != uuid.Nil {
						go kafka.CreateTopic("user_topic_" + dbProfile.ID.String())
					}
					// ==========================================================

					dbAccessRequest, err = as.dao.AccessRequestAccess.CreateAccessRequest(dbProfile.ID, dbEmail.ID)
					if err != nil {
						log.ErrorLogger.Println(err)
						response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
						return &response, err
					}
					bodyPayload := models.ProfileBody{
						ProfileBase:    dbProfile.ProfileBase,
						ProfileId:      dbProfile.Base.ID.String(),
						AccountId:      dbProfile.AccountId,
						Onboarded:      &dbProfile.Onboarded,
						AccessApproved: &dbAccessRequest.Approved,
						Email:          dbEmail.Email,
						TeamId:         dbTeam.Base.ID.String(),
					}

					response := models.SignInResponse{IsSuccessful: true, Code: 2, UserBody: bodyPayload, BearerToken: "bearer_" + dbProfile.Base.ID.String()}
					return &response, nil

				} else {
					response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
					return &response, err
				}
			}
			//  found email but no credential
			// check verified email
			//  if verified create credential and send profile response
			if !dbEmail.Verified {
				dbEmail, err = as.dao.EmailAccess.VerifyEmail(dbEmail)
				if err != nil {
					log.ErrorLogger.Println(err)
					response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
					return &response, err
				}
			}

			_, err = as.dao.CredentialAccess.CreateGoogleCredential(dbEmail.UserId, identifier, authType, providerType)
			if err != nil {
				log.ErrorLogger.Println(err)
				response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
				return &response, err
			}
			// TODO - Update the profile
			dbProfile, err = as.dao.ProfileAccess.GetByUserId(dbEmail.UserId)
			if err != nil {
				log.ErrorLogger.Println(err)
				response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
				return &response, err
			}

			// check primary email
			// var primaryEmail *string
			// if !dbEmail.Primary{
			// 	dbPrimaryEmail, err := as.dao.EmailAccess.FindPrimaryEmailbyUserId(dbEmail.UserId)
			// 	if err != nil {
			// 		log.ErrorLogger.Println(err)

			// 	}
			// 	primaryEmail = &dbPrimaryEmail.Email
			// } else {
			// 	primaryEmail = &dbEmail.Email
			// }

			dbAccessRequest, err = as.dao.AccessRequestAccess.FindAccessRequestByProfileId(dbProfile.ID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					log.ErrorLogger.Println(err)

					dbAccessRequest, err = as.dao.AccessRequestAccess.CreateAccessRequest(dbProfile.ID, dbEmail.ID)
					if err != nil {
						log.ErrorLogger.Println(err)
						response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
						return &response, err
					}
				} else {
					log.ErrorLogger.Println(err)
					response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
					return &response, err
				}
			}
			bodyPayload := models.ProfileBody{
				ProfileBase:    dbProfile.ProfileBase,
				ProfileId:      dbProfile.Base.ID.String(),
				AccountId:      dbProfile.AccountId,
				Onboarded:      &dbProfile.Onboarded,
				AccessApproved: &dbAccessRequest.Approved,
				Email:          dbEmail.Email,
			}

			response := models.SignInResponse{IsSuccessful: true, Code: 2, UserBody: bodyPayload, BearerToken: "bearer_" + dbProfile.Base.ID.String()}
			return &response, nil

		} else {
			log.ErrorLogger.Println(err)
			response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
			return &response, err
		}

	}

	dbProfile, err = as.dao.ProfileAccess.GetByUserId(dbCredential.UserId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
		return &response, err
	}
	dbAccessRequest, err = as.dao.AccessRequestAccess.FindAccessRequestByProfileId(dbProfile.ID)
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	bodyPayload := models.ProfileBody{
		ProfileBase:    dbProfile.ProfileBase,
		ProfileId:      dbProfile.Base.ID.String(),
		AccountId:      dbProfile.AccountId,
		Onboarded:      &dbProfile.Onboarded,
		AccessApproved: &dbAccessRequest.Approved,
		Email:          curEmail,
	}

	//  TODO

	response := models.SignInResponse{IsSuccessful: true, Code: 2, UserBody: bodyPayload, BearerToken: "bearer_" + dbProfile.Base.ID.String()}
	return &response, nil

}
