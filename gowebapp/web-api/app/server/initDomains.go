package server

import (
	"net/http"
	"web-api/app/dummyAuth"
	authController "web-api/app/dummyAuth/controller"
	authDao "web-api/app/dummyAuth/dao"
	authService "web-api/app/dummyAuth/services"
	"web-api/app/middleware"
	"web-api/app/profile"
	profileController "web-api/app/profile/controller"
	profileDao "web-api/app/profile/dao"
	profileService "web-api/app/profile/services"
	"web-api/app/source"
	sourceController "web-api/app/source/controller"
	sourceDao "web-api/app/source/dao"
	sourceService "web-api/app/source/services"
	"web-api/app/utility/respond"
	"web-api/app/view"
	viewController "web-api/app/view/controller"
	viewDao "web-api/app/view/dao"
	viewService "web-api/app/view/services"

	"github.com/go-chi/chi/v5"
)

type Domain struct {
	Auth    *authController.AuthController
	Profile *profileController.ProfileController
	Source  *sourceController.SourceController
	View    *viewController.ViewController
}

func (s *Server) InitDomains() {
	s.initVersion()
	s.initAuth()
	s.initProfile()
	s.initSource()
	// s.initView()
}

func (s *Server) initView() {
	newViewDataAccessOp := viewDao.New(s.DB(), s.cfg)
	newViewService := viewService.NewViewService(*newViewDataAccessOp)
	s.Domain.View = view.RegisterHTTPEndPoints(s.Router, s.validator, newViewService)
}

func (s *Server) initSource() {
	newSourceDataAccessOp := sourceDao.New(s.DB(), s.cfg)
	newSourceService := sourceService.NewSourceService(*newSourceDataAccessOp)
	s.Domain.Source = source.RegisterHTTPEndPoints(s.Router, s.validator, newSourceService)
}

func (s *Server) initProfile() {
	newProfileDataAccessOp := profileDao.New(s.DB(), s.cfg)
	newProfileService := profileService.NewProfileService(*newProfileDataAccessOp)
	s.Domain.Profile = profile.RegisterHTTPEndPoints(s.Router, s.validator, newProfileService)
}

func (s *Server) initAuth() {
	newAuthDataAccessOp := authDao.New(s.DB(), s.cfg)
	newAuthService := authService.NewAuthService(*newAuthDataAccessOp)
	s.Domain.Auth = dummyAuth.RegisterHTTPEndPoints(s.Router, s.validator, newAuthService)
}

func (s *Server) initVersion() {
	s.Router.Route("/", func(router chi.Router) {
		router.Use(middleware.Json)

		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			respond.Json(w, http.StatusOK, map[string]string{"version": s.Version})
		})
	})
}
