/*
* @file internal/server/handlers.go
* @brief this file contains the handlers for the server
 */

package server

import (
	"SmartMeterSystem/cmd/web"
	"net/http"
)

func (s *Server) HomeWebPage(w http.ResponseWriter, r *http.Request) {
	// Log the method and the path
	s.logger.Sugar().Infof("%s\t%s\t%s", s.GetDefaultRouteVersion(), r.URL.Path, r.Method)
	// Set the content type
	w.Header().Set("Content-Type", "text/html")

	// Write the response
	web.HomeWebPage(s.defaultRouteVersion).Render(r.Context(), w)
}

func (s *Server) LoginWebPage(w http.ResponseWriter, r *http.Request) {
	// Log the method and the path
	s.logger.Sugar().Infof("%s\t%s\t%s", s.GetDefaultRouteVersion(), r.URL.Path, r.Method)
	// Set the content type
	w.Header().Set("Content-Type", "text/html")

	// Write the response
	web.LoginWebPage(s.defaultRouteVersion, r.URL.Query().Get("user_type")).Render(r.Context(), w)
}
