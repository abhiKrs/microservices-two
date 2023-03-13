package view

import (
	"web-api/app/view/controller"
	service "web-api/app/view/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	// "gorm.io/gorm"
)

func RegisterHTTPEndPoints(router *chi.Mux, validator *validator.Validate, viewService *service.ViewService) *controller.ViewController {
	// r := chi.NewRouter()

	viewController := controller.NewViewController(viewService, validator)

	// router.Route("/team/{teamId}", func(router chi.Router) {
	// 	router.Get("/", viewController.GetAll)
	// 	router.Post("/", viewController.Create)
	// 	router.Get("/{id}", viewController.GetOne) // POST /posts - Create a new post.
	// 	// router.Put("/{id}", viewController.Update) // POST /posts - Create a new post.
	// 	// router.Delete("/{id}", viewController.Delete)
	// })

	router.Route("/views", func(router chi.Router) {

		router.Get("/", viewController.GetAll)
		router.Post("/", viewController.Create)
		router.Get("/{id}", viewController.GetOne)     // POST /posts - Create a new post.
		router.Put("/{id}", viewController.UpdateView) // POST /posts - Create a new post.
		router.Delete("/{id}", viewController.Delete)

	})
	// router.Route("/view", func(router chi.Router) {
	// 	router.Get("/", viewController.ViewSayHello)
	// })

	return viewController
}
