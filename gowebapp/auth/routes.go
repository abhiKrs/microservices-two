package main

import (
	// "net/http"
	"auth/controller"
	service "auth/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	// "gorm.io/gorm"
)

func RegisterHTTPEndPoints(router *chi.Mux, validator *validator.Validate, authService *service.AuthService) *controller.AuthController {
	// r := chi.NewRouter()

	authController := controller.NewAuthController(authService, validator)

	// router.Route("/api/v1/auth", func(router chi.Router) {
	router.Route("/auth", func(router chi.Router) {
		router.Get("/", authController.SayHello)
		router.Post("/signup", authController.AuthSignup) // POST /posts - Create a new post.
		router.Post("/signin", authController.AuthSignIn) // POST /posts - Create a new post.
		router.Post("/magiclink", authController.AuthMagicLink)
		router.Post("/reset-password-initiate", authController.ResetPasswordInitiate)
		// router.Post("/reset-password-verify", authController.ResetPasswordVerify)
		// router.Post("/reset-password", authController.ResetPassword)
		router.Post("/signout", authController.SignOut)
	})

	return authController
}
