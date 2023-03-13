package controller

import (
	// "errors"
	// "encoding/json"
	"errors"
	"net/http"

	"web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"
	"web-api/app/utility/validate"
)

func (ac *AuthController) ResetPasswordInitiate(w http.ResponseWriter, r *http.Request) {
	log.DebugLogger.Println("started magicLink signin")
	var req models.MagicLinkRequest
	// err := json.NewDecoder(r.Body).Decode(r)
	err := req.Bind(r.Body)
	if err != nil {
		log.ErrorLogger.Println(err)
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(ac.validator, req)
	if errs != nil {
		log.ErrorLogger.Println(errs)
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	res, err := ac.authService.ResetPasswordMagicLink(&req)

	if err != nil {
		log.ErrorLogger.Println(err.Error())
		if errors.Is(err, respond.ErrBadRequest) {
			respond.Error(w, http.StatusBadRequest, err)
			return
		}
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.Json(w, http.StatusCreated, &res)

	log.DebugLogger.Println("finished")

}
