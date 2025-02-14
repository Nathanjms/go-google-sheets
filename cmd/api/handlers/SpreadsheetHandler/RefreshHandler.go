package SpreadsheetHandler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-google-sheets/internal/application"
	"github.com/nathanjms/go-google-sheets/internal/sheets"
)

func RefreshHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		data, err := sheets.FetchSheetData(app.Config, app.Logger) // Use sheets package
		if err != nil {
			app.Logger.Error("Failed to refresh cache", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to refresh cache: "+err.Error())
		}

		app.Cache.Data.Spreadsheet.Data = data
		app.Cache.Data.Spreadsheet.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Cache updated", "data": app.Cache.Data.Spreadsheet.Data})

	}
}
