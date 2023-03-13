package services

import (
	dbConstants "auth/db/constants"
	"auth/model"
	"auth/utility"
	"auth/utility/respond"
	"errors"
	"log"
	"net/http"
)

func (as *AuthService) VerifySigninPasswordAndGetProfile(req *model.SignInRequest) (*model.SignInResponse, error) {
	log.Println("inside pass and get profile")
	var err error

	// get email
	dbEmail, err := as.dao.EmailAccess.FindEmailbyStrAdd(*req.Email)
	if err != nil {
		log.Println("email not found")
		response := model.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
		return &response, err
	}

	// get credential
	dbCredential, err := as.dao.CredentialAccess.GetCredential(dbEmail.UserId, dbConstants.Custom, dbConstants.DBPassword)

	if err != nil {
		log.Println(err)
		if errors.Is(err, respond.ErrNoRecord) {
			log.Println(err)
			response := model.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
			return &response, err
		}
		response := model.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
		return &response, err
	}

	passAlgo := *dbCredential.PasswordAlgo
	log.Println(passAlgo)

	switch passAlgo {
	case "bcrypt":
		err = utility.VerifyPassword(*dbCredential.PasswordHash, req.Credential)

		if err != nil {
			log.Println("wrong password")
			response := model.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{"Password didn't match!!!"}}
			return &response, err
		}
	default:
		response := model.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{http.StatusText(http.StatusInternalServerError)}}
		return &response, respond.ErrInternalServerError
	}

	// get profile
	dbProfile, err := as.dao.ProfileAccess.GetByUserId(dbCredential.UserId)
	if err != nil {
		log.Println("credentials not found")
		response := model.SignInResponse{IsSuccessful: false, Code: 2, Message: []string{err.Error()}}
		return &response, err
	}

	// get Access Request
	dbAccessRequests, err := as.dao.AccessRequestAccess.FindAccessRequestByProfileId(dbProfile.ID)
	if err != nil {
		log.Println(err)
	}

	bodyPayload := model.ProfileBody{
		ProfileBase:    dbProfile.ProfileBase,
		ProfileId:      dbProfile.Base.ID.String(),
		AccountId:      dbProfile.AccountId,
		Onboarded:      &dbProfile.Onboarded,
		AccessApproved: &dbAccessRequests.Approved,
		Email:          dbEmail.Email,
	}

	//  TODO

	response := model.SignInResponse{IsSuccessful: true, Code: 2, UserBody: bodyPayload, BearerToken: "bearer_" + dbCredential.UserId.String()}
	return &response, nil

}
