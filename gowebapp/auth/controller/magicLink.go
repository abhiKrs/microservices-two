package controller

import (
	"auth/model"
	"auth/utility/respond"
	"auth/utility/validate"
	"errors"
	"log"
	"net/http"
)

func (ac *AuthController) AuthMagicLink(w http.ResponseWriter, r *http.Request) {
	log.Println("started magicLink signin")
	var req model.MagicLinkRequest
	// err := json.NewDecoder(r.Body).Decode(r)
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

	res, err := ac.authService.SigninMagicLink(&req)

	if err != nil {
		log.Println(err)
		if errors.Is(err, respond.ErrNoRecord) {
			log.Println("created new user")
			log.Println(err)
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

	log.Println("finished")

}
