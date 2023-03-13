package controller

import (
	// "log"
	"errors"
	"log"
	"net/http"

	"auth/model"
	"auth/utility/respond"
	"auth/utility/validate"
)

// Sigup or register
func (ac *AuthController) AuthSignup(w http.ResponseWriter, r *http.Request) {

	// check request body has proper data type send error if not
	log.Println("started signup")
	var req model.MagicLinkRequest
	err := req.Bind(r.Body)

	if err != nil {
		res := model.MagicLinkResponse{IsSuccessful: false, Email: req.Email, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, res)
		return
	}

	errs := validate.Validate(ac.validator, req)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	res, err := ac.authService.SignupViaMagicLink(&req)

	if err != nil {
		// log.Println("here!!!!")
		log.Println(err.Error())
		if errors.Is(err, respond.ErrBadRequest) {
			log.Println(err.Error())
			respond.Json(w, http.StatusBadRequest, &res)
			return
		} else {
			log.Println(err.Error())
			respond.Json(w, http.StatusBadRequest, &res)
			return
		}
	} else {
		respond.Json(w, http.StatusCreated, &res)
	}

	log.Println("finished")

}
