package controller

import (
	"net/http"

	service "web-api/app/dummyAuth/services"

	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService service.AuthService
	validator   *validator.Validate
}

func NewAuthController(authService *service.AuthService, validator *validator.Validate) *AuthController {
	return &AuthController{
		authService: *authService,
		validator:   validator,
	}
}

func (ac *AuthController) SayHello(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello From Auth!!!"))
}
