package controller

import (
	"errors"
	"net/http"

	"web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"
	"web-api/app/utility/validate"
)

// Sigup or register
func (ac *AuthController) AuthSignup(w http.ResponseWriter, r *http.Request) {

	// check request body has proper data type send error if not
	log.DebugLogger.Println("started signup")
	var req models.MagicLinkRequest
	err := req.Bind(r.Body)

	if err != nil {
		log.ErrorLogger.Println(err)
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

	res, err := ac.authService.SignupViaMagicLink(&req)

	if err != nil {
		// log.Println("here!!!!")
		log.ErrorLogger.Println(err.Error())
		if errors.Is(err, respond.ErrBadRequest) {
			log.ErrorLogger.Println(err.Error())
			respond.Json(w, http.StatusBadRequest, &res)
			return
		} else {
			log.ErrorLogger.Println(err.Error())
			respond.Json(w, http.StatusBadRequest, &res)
			return
		}
	} else {
		respond.Json(w, http.StatusCreated, &res)
	}

	log.DebugLogger.Println("finished")

}
