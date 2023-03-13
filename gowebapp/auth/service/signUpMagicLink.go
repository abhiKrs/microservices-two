package services

import (
	"errors"
	"log"

	// constants "web-api/app/constants"
	dbConstants "auth/db/constants"
	"auth/model"
	"auth/utility/respond"
)

func (as *AuthService) SignupViaMagicLink(req *model.MagicLinkRequest) (*model.MagicLinkResponse, error) {

	// var dbEmail *schema.Email

	dbEmail, err := as.dao.EmailAccess.FindEmailbyStrAdd(req.Email)

	if errors.Is(err, respond.ErrNoRecord) {
		log.Println(err)
		dbEmail, dbMagicLink, err := as.dao.CreateNewUserAndSendMagicLink(req.Email)

		if err != nil {
			log.Println("errorn in transaction")
			res := model.MagicLinkResponse{IsSuccessful: false, Email: req.Email, Message: []string{err.Error()}}
			return &res, err
		}
		// Send Email
		// email.SendEmailviaAws(dbEmail.Email, dbMagicLink.Token.String(), constants.MagicLink, dbConstants.Register,nil, nil)

		responsePayload := model.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime}
		return &responsePayload, nil

	}

	if dbEmail.Verified {
		log.Println(dbEmail)
		log.Println("varified email record found!!!!")
		msg := "already registered user. Sent link to login"
		dbMagicLink, err := as.dao.ResendAuthLink(dbEmail, dbConstants.Login)
		if err != nil {
			responsePayload := model.MagicLinkResponse{IsSuccessful: false, Message: []string{err.Error()}}
			return &responsePayload, err
		}

		responsePayload := model.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime, Message: []string{msg}}
		return &responsePayload, nil

	}

	dbMagicLink, err := as.dao.ResendAuthLink(dbEmail, dbConstants.Register)

	if err != nil {
		responsePayload := model.MagicLinkResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &responsePayload, err
	}

	responsePayload := model.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime}
	return &responsePayload, nil

}
