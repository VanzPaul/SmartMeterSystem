// internal/server/routes/deps.go
package routes

import "go.uber.org/zap"

type ServerDeps interface {
	GetLogger() *zap.Logger
	GetDefaultRouteVersion() string
}
