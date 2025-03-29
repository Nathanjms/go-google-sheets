package sheets

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/joho/godotenv"
	"github.com/nathanjms/go-google-sheets/internal/application"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func FetchSheetData(cfg application.Config, sheetName string, logger *slog.Logger) (application.SpreadsheetData, error) {

	// Load .env file
	err := godotenv.Load() // Loads .env from the current directory
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err) // Handle error appropriately
	}
	ctx := context.Background()

	conf, err := google.CredentialsFromJSON(ctx, []byte(cfg.GoogleServiceAccount), sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		return application.SpreadsheetData{}, fmt.Errorf("could not parse service account credentials: %w", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentials(conf))
	if err != nil {
		return application.SpreadsheetData{}, fmt.Errorf("could not create sheets service: %w", err)
	}

	readRange := sheetName + "!A1:I100" // Loads (up to) 100, so it's the cols that are important
	resp, err := srv.Spreadsheets.Values.Get(cfg.SpreadsheetId, readRange).Do()
	if err != nil {
		return application.SpreadsheetData{}, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	if len(resp.Values) == 0 {
		return application.SpreadsheetData{}, nil
	}

	sheetData := []application.SpreadsheetRow{}
	headers := []string{}
	for i, row := range resp.Values {
		if i == 0 {
			for _, col := range row {
				headers = append(headers, col.(string))
			}
		} else {
			sheetData = append(sheetData, row)
		}
	}

	return application.SpreadsheetData{Headers: headers, Contents: sheetData}, nil
}

func StoreInCache(app *application.Application, sheetName string, data application.SpreadsheetData) {
	// Ensure the main map is initialized
	if app.Cache.Data.Spreadsheets == nil {
		app.Cache.Data.Spreadsheets = make(application.Spreadsheets)
	}

	// Ensure the nested map is initialized
	if _, ok := app.Cache.Data.Spreadsheets[app.Config.SpreadsheetId]; !ok {
		app.Cache.Data.Spreadsheets[app.Config.SpreadsheetId] = make(map[string]application.SpreadsheetCache)
	}

	sheetCache := app.Cache.Data.Spreadsheets[app.Config.SpreadsheetId][sheetName]
	sheetCache.Data = data
	sheetCache.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	app.Cache.Data.Spreadsheets[app.Config.SpreadsheetId][sheetName] = sheetCache
}
