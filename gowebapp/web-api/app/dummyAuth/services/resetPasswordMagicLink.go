package services

import (
	"errors"
	"net/http"

	"web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"

	"gorm.io/gorm"
)

func (as *AuthService) ResetPasswordMagicLink(req *models.MagicLinkRequest) (*models.MagicLinkResponse, error) {

	dbEmail, err := as.dao.EmailAccess.FindEmailbyStrAdd(req.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.ErrorLogger.Println(err)
			return nil, err
		}
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if !dbEmail.Verified {
		return nil, errors.New("invalid request. user not verified yet")
	}
	dbMagicLink, err := as.dao.ResendAuthLink(dbEmail, constants.ResetPassword)

	if err != nil {
		log.ErrorLogger.Println(err)
		log.ErrorLogger.Println("errorn in transaction")
		return nil, err
	}

	responsePayload := models.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime}
	return &responsePayload, nil

}
