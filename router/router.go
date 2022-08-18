package router

import (
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/handlers"
)

// SetupRoutes sets up the routes for the server
func SetupRoutes(e *echo.Echo, h handlers.DBHandler) {
	api := e.Group("/api")
	v1 := api.Group("/v1")
	handlers.RegisterRoutes(v1, h)
}
