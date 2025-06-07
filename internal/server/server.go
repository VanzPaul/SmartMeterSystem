/*
 * @file internal/server/server.go
 * @brief this file contains the server interface
 */

package server

import (
	"SmartMeterSystem/cmd/web"
	"SmartMeterSystem/internal"
	"SmartMeterSystem/internal/server/routes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"go.uber.org/zap"
)

type Role string

const (
	RoleSystemAdmin          Role = "system_admin"
	RoleFinancialAdmin       Role = "financial_admin"
	RoleHRAdmin              Role = "hr_admin"
	RoleCustomerServiceAdmin Role = "customer_service_admin"
	RoleFieldAdmin           Role = "field_admin"
	RoleCashier              Role = "cashier"
	RoleConsumer             Role = "consumer"
)

type ClientType string

const (
	ConsumerType ClientType = "consumer"
	EmployeeType ClientType = "employee"
)

type Server struct {
	port                int
	logger              *zap.Logger
	defaultRouteVersion string
	clienttype          string
	backgroundManager   *BackgroundManager
}

// NewServer creates a new HTTP server instance
func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	defaultRouteVersion := os.Getenv("DEFAULT_ROUTE_VERSION")
	if defaultRouteVersion == "" {
		defaultRouteVersion = "v1"
	}

	logger, loggerErr := internal.NewLogger()
	if loggerErr != nil {
		panic(loggerErr)
	}

	// Create the Server instance
	NewServer := &Server{
		port:                port,
		logger:              logger,
		defaultRouteVersion: defaultRouteVersion,
		clienttype:          "",
	}

	// Initialize background manager
	NewServer.backgroundManager = NewBackgroundManager(NewServer)
	NewServer.backgroundManager.Start()

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Graceful shutdown for background tasks
	server.RegisterOnShutdown(func() {
		NewServer.backgroundManager.Stop()
	})

	return server
}

// Implement ServerDeps interface from routes package
func (s *Server) GetLogger() *zap.Logger {
	return s.logger
}

func (s *Server) GetDefaultRouteVersion() string {
	return s.defaultRouteVersion
}

// RegisterRoutes sets up all HTTP routes with dependencies injected
func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	middlewareGroup := []func(http.Handler) http.Handler{
		s.loggingMiddleware,
	}

	// Create versioned routes with server dependencies injected
	v1Routes := &routes.V1Routes{
		Consumer: routes.V1ConsumerRoute{Deps: s},
		Meter:    routes.V1MeterRoute{Deps: s},
		Employee: routes.V1EmployeeRoute{Deps: s},
	}
	v2Routes := &routes.V2Routes{} // Assuming V2Routes follows similar pattern

	mux.Handle("/", templ.Handler(web.NotFound()))

	// Register v1 routes under /v1/
	mux.Handle("/v1/", http.StripPrefix("/v1", s.applyMiddleware(v1Routes.V1Handler(), middlewareGroup...)))
	// Register v2 routes under /v2/
	mux.Handle("/v2/", http.StripPrefix("/v2", s.applyMiddleware(v2Routes.V2Handler(), middlewareGroup...)))

	// Static files and other routes
	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)
	mux.HandleFunc("/home", s.HomeWebPage)

	return s.corsMiddleware(mux)
}
