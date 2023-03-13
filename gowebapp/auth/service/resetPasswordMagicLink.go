package services

import (
	"errors"
	"log"
	"net/http"

	dbConstants "auth/db/constants"
	"auth/model"
	"auth/utility/respond"
)

func (as *AuthService) ResetPasswordMagicLink(req *model.MagicLinkRequest) (*model.MagicLinkResponse, error) {

	dbEmail, err := as.dao.EmailAccess.FindEmailbyStrAdd(req.Email)

	if err != nil {
		if errors.Is(err, respond.ErrNoRecord) {
			log.Println(err)
			return nil, err
		}
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	if !dbEmail.Verified {
		return nil, errors.New("invalid request. user not verified yet!!!")
	}
	dbMagicLink, err := as.dao.ResendAuthLink(dbEmail, dbConstants.ResetPassword)

	if err != nil {
		log.Println("errorn in transaction")
		return nil, err
	}

	responsePayload := model.MagicLinkResponse{IsSuccessful: true, Email: dbEmail.Email, ExpiryTime: dbMagicLink.ExpiryTime}
	return &responsePayload, nil

}
