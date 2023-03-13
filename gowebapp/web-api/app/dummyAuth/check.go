package dummyAuth

// package main

// import (
// 	"log"
// 	"os"
// 	"regexp"
// 	"time"

// 	"github.com/benc-uk/go-rest-api/pkg/api"
// 	"github.com/benc-uk/go-rest-api/pkg/auth"
// 	"github.com/benc-uk/go-rest-api/pkg/env"
// 	"github.com/benc-uk/go-rest-api/pkg/logging"
// 	"github.com/go-chi/chi"
// 	"github.com/go-chi/chi/middleware"
// )

// type API struct {
// 	*api.Base
// 	service spec.UserService
// }

// var (
// 	healthy     = true               // Simple health flag
// 	version     = "0.2"              // App version number, set at build time with -ldflags "-X 'main.version=1.2.3'"
// 	buildInfo   = "No build details" // Build details, set at build time with -ldflags "-X 'main.buildInfo=Foo bar'"
// 	serviceName = "auth"
// 	defaultPort = 8080
// )

// // Main entry point, will start HTTP service
// func main() {
// 	log.SetOutput(os.Stdout) // Personal preference on log output

// 	// Port to listen on, change the default as you see fit
// 	serverPort := env.GetEnvInt("PORT", defaultPort)

// 	// Use chi for routing
// 	router := chi.NewRouter()

// 	// Our API wraps the Base API and UserService
// 	api := API{
// 		api.NewBase(serviceName, version, buildInfo, healthy),
// 		impl.NewService(serviceName),
// 	}

// 	// Enabling of auth is optional, set via AUTH_CLIENT_ID env var
// 	var validator auth.Validator

// 	// Some basic middleware
// 	router.Use(middleware.RealIP)
// 	router.Use(logging.NewFilteredRequestLogger(regexp.MustCompile(`(^/metrics)|(^/health)`)))
// 	router.Use(middleware.Recoverer)
// 	// Some custom middleware for CORS
// 	router.Use(api.SimpleCORSMiddleware)
// 	// Add Prometheus metrics endpoint, must be before the other routes
// 	api.AddMetricsEndpoint(router, "metrics")

// 	// Add root, health & status middleware
// 	api.AddHealthEndpoint(router, "health")
// 	api.AddStatusEndpoint(router, "status")
// 	api.AddOKEndpoint(router, "")

// 	// Add application routes for this service
// 	api.addRoutes(router, validator)

// 	// Finally start the server
// 	api.StartServer(serverPort, router, 5*time.Second)
// }
