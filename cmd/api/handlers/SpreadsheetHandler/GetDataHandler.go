package SpreadsheetHandler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-google-sheets/internal/application"
	"github.com/nathanjms/go-google-sheets/internal/sheets"
)

func GetDataHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		sheetName := c.QueryParam("sheetName")

		if sheetName == "" {
			sheetName = "Sheet1"
		}

		// Log the name:
		app.Logger.Info("Getting data for sheet sheetName")

		if len(app.Cache.Data.Spreadsheets[app.Config.SpreadsheetId][sheetName].Data.Contents) == 0 ||
			time.Now().UnixNano()/int64(time.Millisecond)-app.Cache.Data.Spreadsheets[app.Config.SpreadsheetId][sheetName].Timestamp > int64(app.Cache.CacheTTL/time.Millisecond) {
			app.Logger.Info("Cache expired. Fetching new data...")
			data, err := sheets.FetchSheetData(app.Config, sheetName, app.Logger) // Use sheets package
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch data: "+err.Error())
			}

			sheets.StoreInCache(app, sheetName, data)
		}

		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "Data Retrieved",
			Data: application.ResponseData{
				"data": app.Cache.Data.Spreadsheets[app.Config.SpreadsheetId][sheetName].Data,
			},
		})
	}
}
