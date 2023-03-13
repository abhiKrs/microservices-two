package services

import (
	"errors"

	"web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"

	"gorm.io/gorm"
)

func (as *AuthService) SigninMagicLink(req *models.MagicLinkRequest) (*models.MagicLinkResponse, error) {
	// Get Email
	dbEmail, err := as.dao.EmailAccess.FindEmailbyStrAdd(req.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.InfoLogger.Println("creating new user")
			log.WarningLogger.Println(err)
			// Create New User
			dbEmail, dbMagicLink, err := as.dao.CreateNewUserAndSendMagicLink(req.Email)

			if err != nil {
				log.ErrorLogger.Println(err)
				log.ErrorLogger.Println("errorn in transaction")
				res := models.MagicLinkResponse{IsSuccessful: false, Email: req.Email, Message: []string{err.Error()}}
				return &res, err
			}
			responsePayload := models.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime, Message: []string{"User not registered. Sent email for signup process."}}
			return &responsePayload, err

		}
		log.ErrorLogger.Println(err)
		res := models.MagicLinkResponse{IsSuccessful: false, Email: req.Email, Message: []string{err.Error()}}
		return &res, err
		// return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if !dbEmail.Verified {
		log.InfoLogger.Println("user not verified yet!!!")
		msg := "Not registered user. Sent link to Register."
		// resend signup link
		dbMagicLink, err := as.dao.ResendAuthLink(dbEmail, constants.Register)
		if err != nil {
			responsePayload := models.MagicLinkResponse{IsSuccessful: false, Message: []string{err.Error()}}
			return &responsePayload, err
		}

		responsePayload := models.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime, Message: []string{msg}}
		return &responsePayload, nil
	}
	dbMagicLink, err := as.dao.ResendAuthLink(dbEmail, constants.Login)

	if err != nil {
		log.ErrorLogger.Println("errorn in transaction")
		res := models.MagicLinkResponse{IsSuccessful: false, Email: req.Email, Message: []string{err.Error()}}
		return &res, err
	}

	responsePayload := models.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime}
	return &responsePayload, nil

}
