package main

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mustafabalila/golang-api/config"
	"github.com/mustafabalila/golang-api/handlers"
	router "github.com/mustafabalila/golang-api/router"
	"github.com/mustafabalila/golang-api/utils/logger"
)

func main() {
	e := echo.New()
	e.Use(logger.LoggerMiddleware())
	e.Use(middleware.Recover())

	cfg := config.GetConfig()
	opt, err := pg.ParseURL(cfg.DatabaseUrl)
	if err != nil {
		e.Logger.Fatal(err)
	}
	db := pg.Connect(opt)

	h := handlers.DBHandler{DB: *db}
	router.SetupRoutes(e, h)
	port := cfg.PORT
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
