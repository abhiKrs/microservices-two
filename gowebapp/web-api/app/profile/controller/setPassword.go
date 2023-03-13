package controller

import (
	"errors"
	"net/http"

	models "web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"
	"web-api/app/utility/validate"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (pc *ProfileController) SetPassword(w http.ResponseWriter, r *http.Request) {

	log.InfoLogger.Println("started password reset")
	profileID := chi.URLParam(r, "id")
	log.DebugLogger.Println(profileID)

	if profileID == "" {
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{"Wrong or missing profile Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}

	pp, err := uuid.Parse(profileID)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}

	var req models.ResetPasswordRequest
	err = req.Bind(r.Body)
	if err != nil {
		log.ErrorLogger.Println(err)
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	log.DebugLogger.Println(req.Password)

	errs := validate.Validate(pc.validator, req)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}
	var res *models.ProfileResponse
	log.DebugLogger.Println("request method is")
	log.DebugLogger.Println(r.Method)

	if req.Password != "" {
		// pc.profileService.
		if r.Method == "POST" {
			res, err = pc.profileService.SetPassword(pp, *&req.Password)
			if err != nil {
				log.ErrorLogger.Println(err)
				respond.Error(w, http.StatusBadRequest, err)
				return
			}
			respond.Json(w, http.StatusAccepted, res)
			return
		} else if r.Method == "PUT" {
			res, err = pc.profileService.ResetPassword(pp, *&req.Password)
			if err != nil {
				log.ErrorLogger.Println(err)
				respond.Error(w, http.StatusBadRequest, err)
				return
			}
			respond.Json(w, http.StatusAccepted, res)
			return
		}

	}
	respond.Error(w, http.StatusBadRequest, errors.New("empty string passed as password"))
	log.DebugLogger.Println("finished")
	return

}
