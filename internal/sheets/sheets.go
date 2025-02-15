package sheets

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathanjms/go-google-sheets/internal/application"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func FetchSheetData(cfg application.Config, logger *slog.Logger) (application.SpreadsheetData, error) {

	// Load .env file
	err := godotenv.Load() // Loads .env from the current directory
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err) // Handle error appropriately
	}
	// Load environment variables
	credentialsJSON := os.Getenv("GOOGLE_SERVICE_ACCOUNT")
	sheetID := os.Getenv("GOOGLE_SHEET_ID")

	ctx := context.Background()

	conf, err := google.CredentialsFromJSON(ctx, []byte(credentialsJSON), sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		return application.SpreadsheetData{}, fmt.Errorf("could not parse service account credentials: %w", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentials(conf))
	if err != nil {
		return application.SpreadsheetData{}, fmt.Errorf("could not create sheets service: %w", err)
	}

	readRange := "Sheet1!A1:I100" // Loads (up to) 100, so it's the cols that are important
	resp, err := srv.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		return application.SpreadsheetData{}, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	if len(resp.Values) == 0 {
		return application.SpreadsheetData{}, nil
	}

	sheetData := []application.SpreadsheetRow{}
	headers := []string{}
	for i, row := range resp.Values {
		fmt.Println(row)
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
