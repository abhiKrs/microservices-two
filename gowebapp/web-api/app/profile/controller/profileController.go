package controller

import (
	"net/http"

	"web-api/app/dummyAuth/models"
	service "web-api/app/profile/services"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"
	"web-api/app/utility/validate"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProfileController struct {
	profileService service.ProfileService
	validator      *validator.Validate
}

func NewProfileController(profileService *service.ProfileService, validator *validator.Validate) *ProfileController {
	return &ProfileController{
		profileService: *profileService,
		validator:      validator,
	}
}

func (pc *ProfileController) ProfileSayHello(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello From Profile!!!"))
}

func (pc *ProfileController) ProfileOnboardSayHello(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello From Profile-Onboarding!!!"))
}

func (pc *ProfileController) GetProfile(w http.ResponseWriter, r *http.Request) {

	log.InfoLogger.Println("getting update")
	if profileID := chi.URLParam(r, "id"); profileID != "" {

		pp, err := uuid.Parse(profileID)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		profile, err := pc.profileService.GetUserProfileById(pp)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		respond.Json(w, http.StatusAccepted, profile)
		return

	} else {
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{"Wrong or missing profile Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	// log.Println("done")
}

func (pc *ProfileController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log.InfoLogger.Println("started profile update")
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

	// if req.Password != nil {
	// 	// pc.profileService.
	// 	_, err := pc.profileService.SetResetPassword(pp, *req.Password)
	// 	if err != nil {
	// 		log.ErrorLogger.Println(err)
	// 	}
	// }

	profile, err := pc.profileService.UpdateProfile(pp, &req.ProfileBase)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	respond.Json(w, http.StatusAccepted, profile)
	return
}
