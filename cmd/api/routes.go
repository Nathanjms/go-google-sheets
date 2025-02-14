package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-google-sheets/cmd/api/handlers/SpreadsheetHandler"
	"github.com/nathanjms/go-google-sheets/internal/application"
)

func InitRoutes(e *echo.Echo, app *application.Application) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "Hello, World!",
		})
	})
	e.GET("status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "OK",
		})
	})

	// Spreadsheet endpoints
	e.GET("data", SpreadsheetHandler.GetDataHandler(app))
	e.POST("reload", SpreadsheetHandler.RefreshHandler(app))
}
