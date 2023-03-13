package services

import (
	"errors"

	dbConstants "web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"

	// "web-api/app/utility/email"

	"gorm.io/gorm"
)

func (as *AuthService) SignupViaMagicLink(req *models.MagicLinkRequest) (*models.MagicLinkResponse, error) {

	// var dbEmail *schema.Email

	dbEmail, err := as.dao.EmailAccess.FindEmailbyStrAdd(req.Email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.ErrorLogger.Println(err)
		dbEmail, dbMagicLink, err := as.dao.CreateNewUserAndSendMagicLink(req.Email)

		if err != nil {
			log.ErrorLogger.Println("errorn in transaction")
			res := models.MagicLinkResponse{IsSuccessful: false, Email: req.Email, Message: []string{err.Error()}}
			return &res, err
		}
		// Send Email
		// email.SendEmailviaAws(dbEmail.Email, dbMagicLink.Token.String(), constants.MagicLink, dbConstants.Register,nil, nil)

		responsePayload := models.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime}
		return &responsePayload, nil

	}

	if dbEmail.Verified {
		log.DebugLogger.Println(dbEmail)
		log.InfoLogger.Println("varified email record found!!!!")
		msg := "already registered user. Sent link to login"
		dbMagicLink, err := as.dao.ResendAuthLink(dbEmail, dbConstants.Login)
		if err != nil {
			responsePayload := models.MagicLinkResponse{IsSuccessful: false, Message: []string{err.Error()}}
			return &responsePayload, err
		}

		responsePayload := models.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime, Message: []string{msg}}
		return &responsePayload, nil

	}

	dbMagicLink, err := as.dao.ResendAuthLink(dbEmail, dbConstants.Register)

	if err != nil {
		responsePayload := models.MagicLinkResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &responsePayload, err
	}

	responsePayload := models.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime}
	return &responsePayload, nil

}
