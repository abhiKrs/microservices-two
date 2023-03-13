package source

import (
	"web-api/app/source/controller"
	service "web-api/app/source/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func RegisterHTTPEndPoints(router *chi.Mux, validator *validator.Validate, sourceService *service.SourceService) *controller.SourceController {

	sourceController := controller.NewSourceController(sourceService, validator)

	router.Route("/source", func(router chi.Router) {

		router.Get("/", sourceController.GetAll)
		router.Post("/", sourceController.Create)
		router.Get("/{id}", sourceController.GetOne)
		router.Put("/{id}", sourceController.UpdateSource)
		router.Delete("/{id}", sourceController.Delete)

	})

	router.Route("/logfire.sh", func(router chi.Router) {
		router.Post("/", sourceController.SendLog)
	})

	return sourceController
}
