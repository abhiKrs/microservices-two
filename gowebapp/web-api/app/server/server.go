package server

import (
	"context"
	"net/http"
	"time"
	"web-api/app/config"

	pgsql "web-api/app/database"
	customMiddleware "web-api/app/middleware"
	log "web-api/app/utility/logger"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Server struct {
	Version string
	cfg     *config.Config

	db *gorm.DB

	validator  *validator.Validate
	Router     *chi.Mux
	httpServer *http.Server

	Domain
}

type Options func(opts *Server) error

func defaultServer() *Server {
	return &Server{
		cfg:    config.New(),
		Router: chi.NewRouter(),
	}
}

func New(opts ...Options) *Server {
	s := defaultServer()

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			log.ErrorLogger.Fatalln(err)
		}
	}
	return s
}

func (s *Server) Init(version string) {
	s.Version = version
	// log.Println(*s.cfg)
	s.newDatabase()
	s.newValidator()
	s.newRouter()
	s.setGlobalMiddleware()
	s.InitDomains()
}

func (s *Server) newDatabase() {
	if s.cfg.PGDatabase.DRIVER == "" {
		log.ErrorLogger.Fatal("please fill in database credentials in .env file or set in environment variable")
	}
	// pgsql.Connect(s.cfg.PGDatabase)
	// s.db = pgsql.DB

	s.db = pgsql.Connect(s.cfg.PGDatabase)

}

func (s *Server) newValidator() {
	s.validator = validator.New()
}

func (s *Server) newRouter() {
	s.Router = chi.NewRouter()
}

func (s *Server) setGlobalMiddleware() {
	s.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "endpoint not found"}`))
	})
	s.Router.Use(customMiddleware.Json)
	s.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "http://localhost"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	if s.cfg.Api.REQUEST_LOG {
		s.Router.Use(chiMiddleware.Logger)
	}
	s.Router.Use(chiMiddleware.Recoverer)

}

func (s *Server) Run() {
	log.InfoLogger.Printf("Database hosted at: %s, database name:%s\n", s.cfg.PGDatabase.HOST, s.cfg.PGDatabase.DB)

	start(s)

	_ = gracefulShutdown(context.Background(), s)
}

func (s *Server) Config() *config.Config {
	return s.cfg
}

func (s *Server) DB() *gorm.DB {
	return s.db
}

func start(s *Server) {
	log.InfoLogger.Printf("Serving at %s:%s\n", s.cfg.Api.HOST, s.cfg.Api.PORT)
	log.DebugLogger.Printf("Inside %s:Environment\n", s.cfg.Api.ENV)
	err := http.ListenAndServe(":"+s.cfg.Api.PORT, s.Router)
	if err != nil {
		log.ErrorLogger.Println("error in serving the code")
		log.ErrorLogger.Fatal(err)
	}
}

func gracefulShutdown(ctx context.Context, s *Server) error {

	log.InfoLogger.Println("Shutting down...")

	ctx, shutdown := context.WithTimeout(ctx, s.Config().Api.GracefulTimeout*time.Second)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}
