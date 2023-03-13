package middleware

import (
	"net/http"

	"github.com/golang/gddo/httputil/header"
)

func Json(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "" {
			value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
			if value != "application/json" {
				msg := "Content-Type header is not application/json"
				http.Error(w, msg, http.StatusUnsupportedMediaType)
				return
			}
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
