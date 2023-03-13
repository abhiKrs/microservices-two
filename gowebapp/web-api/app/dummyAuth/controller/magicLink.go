package controller

import (
	"errors"
	"net/http"

	"web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"
	"web-api/app/utility/validate"

	"gorm.io/gorm"
)

func (ac *AuthController) AuthMagicLink(w http.ResponseWriter, r *http.Request) {
	log.DebugLogger.Println("started magicLink signin")
	var req models.MagicLinkRequest
	// err := json.NewDecoder(r.Body).Decode(r)
	err := req.Bind(r.Body)
	if err != nil {
		res := models.MagicLinkResponse{IsSuccessful: false, Email: req.Email, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, res)
		return
	}

	errs := validate.Validate(ac.validator, req)
	if errs != nil {
		log.ErrorLogger.Println(errs)
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	res, err := ac.authService.SigninMagicLink(&req)

	if err != nil {
		log.ErrorLogger.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.InfoLogger.Println("created new user")
			log.ErrorLogger.Println(err)
			// Create New User
			respond.Json(w, http.StatusCreated, &res)
			return
		} else if errors.Is(err, respond.ErrBadRequest) {
			respond.Json(w, http.StatusBadRequest, &res)
			return
		} else {
			respond.Json(w, http.StatusBadRequest, &res)
			return
		}
	}
	respond.Json(w, http.StatusCreated, &res)

	log.DebugLogger.Println("finished")

}
