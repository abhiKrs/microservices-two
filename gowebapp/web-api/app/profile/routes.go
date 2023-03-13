package profile

import (
	// "net/http"
	"web-api/app/profile/controller"
	service "web-api/app/profile/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	// "gorm.io/gorm"
)

func RegisterHTTPEndPoints(router *chi.Mux, validator *validator.Validate, profileService *service.ProfileService) *controller.ProfileController {
	// r := chi.NewRouter()

	profileController := controller.NewProfileController(profileService, validator)

	router.Route("/profile", func(router chi.Router) {
		router.Get("/", profileController.ProfileSayHello)
		router.Get("/{id}", profileController.GetProfile)    // POST /posts - Create a new post.
		router.Put("/{id}", profileController.UpdateProfile) // POST /posts - Create a new post.

	})

	router.Route("/onboard", func(router chi.Router) {
		router.Get("/", profileController.ProfileOnboardSayHello)
		router.Put("/{id}", profileController.OnboardProfile) // POST /posts - Create a new post.

	})
	router.Route("/set-password", func(router chi.Router) {
		router.Post("/{id}", profileController.SetPassword) // POST /posts - Create a new post.
		router.Put("/{id}", profileController.SetPassword)  // POST /posts - Create a new post.
	})

	return profileController
}
