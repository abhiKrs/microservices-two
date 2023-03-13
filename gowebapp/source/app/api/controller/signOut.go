package controller

import (
	"net/http"
	"strings"

	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"
)

type SignOut struct {
	IsSuccessful bool   `json:"isSuccessful"`
	Msg          string `json:"msg"`
}

func (ac *AuthController) SignOut(w http.ResponseWriter, r *http.Request) {

	log.DebugLogger.Println("started signout")
	// var req models.SignInRequest
	token := r.Header.Get("Authorization")
	tokenVal := strings.Split(token, "_")[1]

	payload := SignOut{IsSuccessful: true, Msg: "dummy test to logout" + tokenVal}

	w.WriteHeader(http.StatusAccepted)
	// w.Write(json.Marshal(payload))
	respond.Json(w, http.StatusAccepted, payload)

}
