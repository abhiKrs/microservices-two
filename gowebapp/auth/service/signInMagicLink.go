package services

import (
	dbConstants "auth/db/constants"
	"auth/model"
	"auth/utility/respond"
	"errors"
	"log"
)

func (as *AuthService) SigninMagicLink(req *model.MagicLinkRequest) (*model.MagicLinkResponse, error) {
	// Get Email
	dbEmail, err := as.dao.EmailAccess.FindEmailbyStrAdd(req.Email)

	if err != nil {
		if errors.Is(err, respond.ErrNoRecord) {
			log.Println("creating new user")
			log.Println(err)
			// Create New User
			dbEmail, dbMagicLink, err := as.dao.CreateNewUserAndSendMagicLink(req.Email)

			if err != nil {
				log.Println("errorn in transaction")
				res := model.MagicLinkResponse{IsSuccessful: false, Email: req.Email, Message: []string{err.Error()}}
				return &res, err
			}
			responsePayload := model.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime, Message: []string{"User not registered. Sent email for signup process."}}
			return &responsePayload, err

		}
		log.Println(err)
		res := model.MagicLinkResponse{IsSuccessful: false, Email: req.Email, Message: []string{err.Error()}}
		return &res, err
		// return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if !dbEmail.Verified {
		log.Println("user not verified yet!!!")
		msg := "Not registered user. Sent link to Register."
		// resend signup link
		dbMagicLink, err := as.dao.ResendAuthLink(dbEmail, dbConstants.Register)
		if err != nil {
			responsePayload := model.MagicLinkResponse{IsSuccessful: false, Message: []string{err.Error()}}
			return &responsePayload, err
		}

		responsePayload := model.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime, Message: []string{msg}}
		return &responsePayload, nil
	}
	dbMagicLink, err := as.dao.ResendAuthLink(dbEmail, dbConstants.Login)

	if err != nil {
		log.Println("errorn in transaction")
		res := model.MagicLinkResponse{IsSuccessful: false, Email: req.Email, Message: []string{err.Error()}}
		return &res, err
	}

	responsePayload := model.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime}
	return &responsePayload, nil

}
