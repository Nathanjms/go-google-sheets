package main

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nathanjms/go-google-sheets/internal/application"
	"github.com/nathanjms/go-google-sheets/internal/env"
)

func serveHttp(app *application.Application) error {
	e := echo.New()

	originsFromEnv := env.GetString("ALLOWED_ORIGINS_BY_COMMA", "http://localhost:3000")
	allowedOrigins := []string{}

	for _, origin := range strings.Split(originsFromEnv, ",") {
		allowedOrigins = append(allowedOrigins, strings.Trim(origin, " "))
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentEncoding},
		AllowCredentials: true,
	}))
	InitRoutes(e, app)

	app.Logger.Info("Starting server on port " + strconv.Itoa(app.Config.HTTPPort))

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(app.Config.HTTPPort)))
	return nil
}
