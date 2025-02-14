package main

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nathanjms/go-google-sheets/internal/application"
)

func serveHttp(app *application.Application) error {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://*.nathanjms.co.uk"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentEncoding},
		AllowCredentials: true,
	}))
	InitRoutes(e, app)

	app.Logger.Info("Starting server on port " + strconv.Itoa(app.Config.HTTPPort))

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(app.Config.HTTPPort)))
	return nil
}
