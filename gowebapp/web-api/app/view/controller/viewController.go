package controller

import (
	"errors"
	"net/http"
	"strings"

	"web-api/app/view/models"
	service "web-api/app/view/services"

	// "web-api/app/utility/kafka"
	log "web-api/app/utility/logger"
	// myRedis "web-api/app/utility/redis"
	"web-api/app/utility/respond"
	"web-api/app/utility/validate"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	// "github.com/google/uuid"
)

type ViewController struct {
	viewService service.ViewService
	validator   *validator.Validate
}

func NewViewController(viewService *service.ViewService, validator *validator.Validate) *ViewController {
	return &ViewController{
		viewService: *viewService,
		validator:   validator,
	}
}

func (pc *ViewController) ViewSayHello(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello From View!!!"))
}

// func (pc *ViewController) ViewOnboardSayHello(w http.ResponseWriter, r *http.Request) {

// 	w.Write([]byte("Hello From View-Onboarding!!!"))
// }

func (pc *ViewController) GetAll(w http.ResponseWriter, r *http.Request) {

	log.InfoLogger.Println("getting all views")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token
	splitToken := strings.Split(reqToken, "bearer_")
	// log.Println(l)
	// x := strings.Join(l[1:], "")
	// splitToken := strings.Split(reqToken, "bearer_")
	// To refactor after jwt
	log.DebugLogger.Println(splitToken)
	profileId := strings.Join(splitToken[1:], "")
	log.InfoLogger.Println(profileId)
	if profileId != "" {

		pId, err := uuid.Parse(profileId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetAllViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		payload, err := pc.viewService.GetViewsByProfileId(pId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetAllViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		respond.Json(w, http.StatusAccepted, payload)
		return

	} else {
		response := models.GetAllViewResponse{IsSuccessful: false, Message: []string{"Wrong or missing source Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	// log.Println("done")
}

func (pc *ViewController) Create(w http.ResponseWriter, r *http.Request) {

	log.InfoLogger.Println("Creating a View")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token

	splitToken := strings.Split(reqToken, "bearer_")
	log.DebugLogger.Println(splitToken)
	profileId := strings.Join(splitToken[1:], "")
	log.InfoLogger.Println(profileId)

	var reqBody models.CreateViewRequest
	err := reqBody.Bind(r.Body)
	if err != nil {
		log.ErrorLogger.Println(err)
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	log.DebugLogger.Println(reqBody)

	errs := validate.Validate(pc.validator, reqBody)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	if profileId != "" {
		pId, err := uuid.Parse(profileId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.CreateViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		payload, err := pc.viewService.CreateView(reqBody, pId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.CreateViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}

		respond.Json(w, http.StatusCreated, payload)
		return

	} else {
		response := models.CreateViewResponse{IsSuccessful: false, Message: []string{"Wrong or missing source Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	// log.Println("done")
}

func (pc *ViewController) GetOne(w http.ResponseWriter, r *http.Request) {

	log.InfoLogger.Println("Creating a View")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token

	splitToken := strings.Split(reqToken, "bearer_")
	// log.DebugLogger.Println(splitToken)
	profileId := strings.Join(splitToken[1:], "")
	// log.InfoLogger.Println(profileId)
	pId, err := uuid.Parse(profileId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.GetViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}

	// log.InfoLogger.Println(profileId)
	if viewID := chi.URLParam(r, "id"); viewID != "" {
		vId, err := uuid.Parse(viewID)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		payload, err := pc.viewService.GetViewById(vId, pId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		respond.Json(w, http.StatusAccepted, payload)
		return

	} else {
		response := models.GetViewResponse{IsSuccessful: false, Message: []string{"Wrong or missing source Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	// log.Println("done")
}

func (pc *ViewController) Delete(w http.ResponseWriter, r *http.Request) {

	log.InfoLogger.Println("Deleting a View")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token

	splitToken := strings.Split(reqToken, "bearer_")
	log.DebugLogger.Println(splitToken)
	profileId := strings.Join(splitToken[1:], "")
	log.InfoLogger.Println(profileId)
	pId, err := uuid.Parse(profileId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.DeleteViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	log.InfoLogger.Println(profileId)
	if viewID := chi.URLParam(r, "id"); viewID != "" {
		vId, err := uuid.Parse(viewID)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.DeleteViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		payload, err := pc.viewService.Delete(vId, pId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.DeleteViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		respond.Json(w, http.StatusAccepted, payload)
		return

	} else {
		response := models.DeleteViewResponse{IsSuccessful: false, Message: []string{"Wrong or missing source Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	// log.Println("done")
}

func (pc *ViewController) UpdateView(w http.ResponseWriter, r *http.Request) {
	log.InfoLogger.Println("Updating a View")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token

	splitToken := strings.Split(reqToken, "bearer_")
	log.DebugLogger.Println(splitToken)
	profileId := strings.Join(splitToken[1:], "")
	log.InfoLogger.Println(profileId)
	pId, err := uuid.Parse(profileId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.GetViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}

	var reqBody models.UpdateViewRequest
	err = reqBody.Bind(r.Body)
	if err != nil {
		log.ErrorLogger.Println(err)
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(pc.validator, reqBody)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	if viewID := chi.URLParam(r, "id"); viewID != "" {
		vId, err := uuid.Parse(viewID)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		payload, err := pc.viewService.UpdateView(pId, vId, &reqBody)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		respond.Json(w, http.StatusAccepted, payload)
		return

	} else {
		response := models.GetViewResponse{IsSuccessful: false, Message: []string{"Wrong or missing source Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
}
