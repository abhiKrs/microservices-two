package controller

import (
	// "errors"
	// "encoding/json"
	"auth/model"
	"auth/utility/respond"
	"auth/utility/validate"
	"errors"
	"log"
	"net/http"
)

func (ac *AuthController) ResetPasswordInitiate(w http.ResponseWriter, r *http.Request) {
	log.Println("started magicLink signin")
	var req model.MagicLinkRequest
	// err := json.NewDecoder(r.Body).Decode(r)
	err := req.Bind(r.Body)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(ac.validator, req)
	if errs != nil {

		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	res, err := ac.authService.ResetPasswordMagicLink(&req)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, respond.ErrBadRequest) {
			respond.Error(w, http.StatusBadRequest, err)
			return
		}
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.Json(w, http.StatusCreated, &res)

	log.Println("finished")

}
