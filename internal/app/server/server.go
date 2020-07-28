package server

import (
	"github.com/labstack/echo/v4"
)

// Server ...
type Server struct {
	// config *Config
	router *echo.Echo
	// store  *store.Store
}
