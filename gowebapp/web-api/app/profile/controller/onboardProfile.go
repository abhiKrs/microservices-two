package controller

import (
	"net/http"

	"web-api/app/dummyAuth/models"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"
	"web-api/app/utility/validate"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (pc *ProfileController) OnboardProfile(w http.ResponseWriter, r *http.Request) {
	log.InfoLogger.Println("started onboarding")
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

	var req models.ProfileUpdateRequest
	err = req.Bind(r.Body)
	if err != nil {
		log.ErrorLogger.Println(err)
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(pc.validator, req)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	profile, err := pc.profileService.OnboardProfileAndPassword(pp, &req)
	if err != nil {
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	respond.Json(w, http.StatusAccepted, profile)
	return

}
