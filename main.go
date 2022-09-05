package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mustafabalila/golang-api/config"
	"github.com/mustafabalila/golang-api/router"
	"github.com/mustafabalila/golang-api/utils/logger"
)

func main() {
	e := echo.New()
	e.Use(logger.LoggerMiddleware())
	e.Use(middleware.Recover())

	cfg := config.GetConfig()

	router.SetupRoutes(e)
	port := cfg.PORT
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
