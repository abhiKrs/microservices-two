package controller

import (
	"errors"
	"net/http"

	constants "web-api/app/constants"
	models "web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"
	"web-api/app/utility/validate"

	"gorm.io/gorm"
)

func (ac *AuthController) AuthSignIn(w http.ResponseWriter, r *http.Request) {

	log.InfoLogger.Println("started signin")
	var req models.SignInRequest
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

	switch req.AuthType {
	case uint8(constants.MagicLink.EnumIndex()):
		log.DebugLogger.Println("inside" + constants.MagicLink.String())
		log.DebugLogger.Println(req.AuthType)
		res, err := ac.authService.VerifyMagicLink(&req)
		if err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				respond.Json(w, http.StatusNotFound, &res)
				return
			} else if errors.Is(err, respond.ErrBadRequest) {
				respond.Json(w, http.StatusBadRequest, &res)
				return
			} else {
				log.ErrorLogger.Println(err.Error())
				respond.Json(w, http.StatusBadRequest, &res)
				return
			}
			// respond.Json(w, http.StatusNotFound, res)
		}
		respond.Json(w, http.StatusOK, &res)
		return

	case uint8(constants.Password.EnumIndex()):
		log.InfoLogger.Println("inside password signin")
		if &req.Email == nil || &req.Credential == nil {
			http.Error(w, "missing credentials for authtype password", http.StatusBadRequest)
			return
		}

		res, err := ac.authService.VerifySigninPasswordAndGetProfile(&req)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				respond.Json(w, http.StatusNotFound, &res)
				return
			} else if errors.Is(err, respond.ErrBadRequest) {
				respond.Json(w, http.StatusBadRequest, &res)
				return
			} else {
				log.ErrorLogger.Println(err.Error())
				respond.Json(w, http.StatusBadRequest, &res)
				return
			}
		}
		respond.Json(w, http.StatusOK, res)
		return

	case uint8(constants.Google.EnumIndex()):
		log.InfoLogger.Println("inside google signin")
		// if &req.Email == nil || &req.Credential == nil {
		// 	http.Error(w, "missing credentials for authtype password", http.StatusBadRequest)
		// 	return
		// }

		res, err := ac.authService.VerifyGoogleSigninAndGetProfile(&req)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.ErrorLogger.Println(err)
				respond.Json(w, http.StatusNotFound, &res)
				return
			} else if errors.Is(err, respond.ErrBadRequest) {
				log.ErrorLogger.Println(err)
				respond.Json(w, http.StatusBadRequest, &res)
				return
			} else {
				log.ErrorLogger.Println(err)
				respond.Json(w, http.StatusBadRequest, &res)
				return
			}
		}
		respond.Json(w, http.StatusOK, res)
		return

	default:
		log.WarningLogger.Print("wrong auth type")
		http.Error(w, "Not Implemented", http.StatusBadRequest)
		// return http.StatusText(http.StatusInternalServerError), http.StatusRequestEntityTooLarge
		return
	}

}
