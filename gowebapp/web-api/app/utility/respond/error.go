package respond

import (
	"encoding/json"
	"errors"
	"net/http"

	log "web-api/app/utility/logger"
)

type ErrorsMessage struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message"`
}

// type ErrorMessage struct {
// 	IsSuccessful bool     `json:"isSuccessful"`
// 	Message      string `json:"message"`
// }

var (
	ErrBadRequest          = errors.New("bad request")
	ErrNoRecord            = errors.New("no record found")
	ErrInternalServerError = errors.New("internal server error")

	ErrDatabase       = errors.New("connecting to database")
	ErrInvalidRequest = errors.New("invalid request")
	ErrTokenExpired   = errors.New("token expired")
)

func Errors(w http.ResponseWriter, statusCode int, errors []string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)

	if errors == nil {
		write(w, nil)
		return
	}

	p := ErrorsMessage{
		IsSuccessful: false,
		Message:      errors,
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	if string(data) == "null" {
		return
	}

	write(w, data)
}

func Error(w http.ResponseWriter, statusCode int, message error) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)

	// var p ErrorsMessage
	if message == nil {
		write(w, nil)
		return
	}

	p := ErrorsMessage{
		IsSuccessful: false,
		Message:      []string{message.Error()},
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	if string(data) == "null" {
		return
	}

	write(w, data)
}

func write(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
}
