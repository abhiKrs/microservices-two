package api

import (
	"source/app/api/controller"
	service "source/app/api/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	// "gorm.io/gorm"
)

func RegisterHTTPEndPoints(router *chi.Mux, validator *validator.Validate, sourceService *service.SourceService) *controller.sourceController {
	// r := chi.NewRouter()

	sourceController := controller.NewSourceController(sourceService, validator)

	// router.Route("/api/v1/Source", func(router chi.Router) {
	router.Route("/source", func(router chi.Router) {
		router.Get("/", sourceController.GetAll)
		router.Post("/", sourceController.Create)
		router.Get("/{id}", sourceController.GetOne) // POST /posts - Create a new post.
		router.Put("/{id}", sourceController.Update) // POST /posts - Create a new post.
		router.Delete("/{id}", sourceController.Delete)

	})

	return sourceController
}
