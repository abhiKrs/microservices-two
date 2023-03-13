package server

import (
	"net/http"
	"source/app/api"
	sourceController "source/app/api/controller"
	sourceDao "source/app/api/dao"
	sourceService "source/app/api/services"
	"source/app/middleware"
	"source/app/utility/respond"

	"github.com/go-chi/chi/v5"
)

type Domain struct {
	Source *sourceController.SourceController
}

func (s *Server) InitDomains() {
	s.initVersion()
	s.initSource()
}

func (s *Server) initVersion() {
	s.Router.Route("/", func(router chi.Router) {
		router.Use(middleware.Json)

		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			respond.Json(w, http.StatusOK, map[string]string{"version": s.Version})
		})
	})
}

func (s *Server) initSource() {
	newSourceDataAccessOp := sourceDao.New(s.DB(), s.cfg)
	newSourceService := sourceService.NewSourceService(*newSourceDataAccessOp)
	s.Domain.Source = api.RegisterHTTPEndPoints(s.Router, s.validator, newSourceService)
}
