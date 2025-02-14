package SpreadsheetHandler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-google-sheets/internal/application"
	"github.com/nathanjms/go-google-sheets/internal/cache"
	"github.com/nathanjms/go-google-sheets/internal/sheets"
)

func GetDataHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		now := time.Now().UnixNano() / int64(time.Millisecond)

		if len(cache.Spreadsheet.Data) == 0 || now-cache.Spreadsheet.Timestamp > int64(cache.CacheInstance.CacheTTL/time.Millisecond) {
			app.Logger.Info("Cache expired. Fetching new data...")
			data, err := sheets.FetchSheetData(app.Config, app.Logger) // Use sheets package
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch data: "+err.Error())
			}
			cache.Spreadsheet.Data = data
			cache.Spreadsheet.Timestamp = now
		}

		return c.JSON(http.StatusOK, cache.Spreadsheet.Data)
	}
}
