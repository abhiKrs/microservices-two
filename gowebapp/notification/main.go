package main

import (
	"encoding/json"
	"fmt"
	"io"

	// "log"
	"net/http"

	"notification/app/config"
	"notification/app/constants"
	"notification/app/email"
	log "notification/app/utility"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type ErrorsMessage struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message"`
}

type EmailResponseBody struct {
	IsSuccessful bool `json:"isSuccessful"`
}

type EmailRequestBody struct {
	Email       string                     `json:"email"`
	Token       *string                    `json:"token,omitempty"`
	EmailType   constants.EmailType        `json:"emailType"`
	LinkType    *constants.DBMagicLinkType `json:"linkType,omitempty"`
	SuccessType constants.SuccessType      `json:"successType"`
}

func (r *EmailRequestBody) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

func Error(w http.ResponseWriter, statusCode int, message error) {
	w.WriteHeader(statusCode)
	p := ErrorsMessage{
		IsSuccessful: false,
		Message:      []string{message.Error()},
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	_, err = w.Write(data)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
}

func Json(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.ErrorLogger.Println(err)
		Error(w, http.StatusInternalServerError, err)
		return
	}

	if string(data) == "null" {
		_, _ = w.Write([]byte("[]"))
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.ErrorLogger.Println(err)
		Error(w, http.StatusInternalServerError, err)
		return
	}
}

func sendEmail(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/problem+json")

	fmt.Println("Processing email...")
	var req EmailRequestBody
	err := req.Bind(r.Body)
	if err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	// err = email.SendEmailviaAws(req.Email, req.Token, req.EmailType, req.LinkType, &req.SuccessType, config.CFG)
	err = email.SendEmailViaMailGun(req.Email, req.Token, req.EmailType, req.LinkType, &req.SuccessType, config.CFG)
	if err != nil {
		log.ErrorLogger.Println(err)
		Error(w, http.StatusBadRequest, err)
		return
	}

	fmt.Println("Finished processing email")
	Json(w, http.StatusOK, EmailResponseBody{IsSuccessful: true})

}

func main() {
	var appPort = "6002"

	log.DebugLogger.Println(appPort)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)

	r.Post("/email", sendEmail)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("welcome"))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("welcome, imok"))
	})

	err := http.ListenAndServe(":"+appPort, r)

	if err != nil {
		log.ErrorLogger.Fatal(err)
	}
}
