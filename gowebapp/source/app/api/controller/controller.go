package controller

import (
	"net/http"

	service "source/app/api/services"

	"github.com/go-playground/validator/v10"
)

type SourceController struct {
	authService service.AuthService
	validator   *validator.Validate
}

func NewSourceController(authService *service.AuthService, validator *validator.Validate) *SourceController {
	return &SourceController{
		authService: *authService,
		validator:   validator,
	}
}

func (ac *SourceController) SayHello(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello From Source!!!"))
}
