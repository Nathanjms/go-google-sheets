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

func FetchSheetData(cfg application.Config, logger *slog.Logger) ([]application.SpreadsheetRow, error) {

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
		return nil, fmt.Errorf("could not parse service account credentials: %w", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentials(conf))
	if err != nil {
		return nil, fmt.Errorf("could not create sheets service: %w", err)
	}

	readRange := "Sheet1!A1:A4" //  Adjust range to match the expected data structure
	resp, err := srv.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	if len(resp.Values) == 0 {
		return []application.SpreadsheetRow{}, nil
	}

	var sheetData []application.SpreadsheetRow
	for _, row := range resp.Values {
		data := application.SpreadsheetRow{
			Name: fmt.Sprint(row[0]),
		}
		sheetData = append(sheetData, data)
	}

	return sheetData, nil
}
