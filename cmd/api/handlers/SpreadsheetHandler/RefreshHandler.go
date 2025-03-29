package SpreadsheetHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-google-sheets/internal/application"
	"github.com/nathanjms/go-google-sheets/internal/sheets"
)

func RefreshHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		sheetName := c.QueryParam("sheetName")

		if sheetName == "" {
			sheetName = "Sheet1"
		}
		data, err := sheets.FetchSheetData(app.Config, sheetName, app.Logger) // Use sheets package
		if err != nil {
			app.Logger.Error("Failed to refresh cache", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to refresh cache: "+err.Error())
		}

		sheets.StoreInCache(app, sheetName, data)

		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "Cache updated",
			Data:    application.ResponseData{"data": app.Cache.Data.Spreadsheets[app.Config.SpreadsheetId][sheetName].Data},
		})

	}
}
