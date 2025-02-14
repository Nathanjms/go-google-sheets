package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetRow struct {
	Name string `json:"name"`
}

type Cache struct {
	Spreadsheet struct {
		Data      []SpreadsheetRow `json:"data"`
		Timestamp int64            `json:"timestamp"`
	} `json:"spreadsheet"`
}

var cache Cache

const (
	cacheTTL = 60 * 60 * time.Second // Cache duration: 60 minutes
)

func main() {
	// Load .env file
	err := godotenv.Load() // Loads .env from the current directory
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err) // Handle error appropriately
	}
	// Load environment variables
	serviceAccountJSON := os.Getenv("GOOGLE_SERVICE_ACCOUNT")
	sheetID := os.Getenv("GOOGLE_SHEET_ID")

	if serviceAccountJSON == "" || sheetID == "" {
		log.Fatal("GOOGLE_SERVICE_ACCOUNT and GOOGLE_SHEET_ID environment variables must be set")
	}

	// Initialize cache
	cache = Cache{}

	// Set up HTTP server using chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/refresh", func(w http.ResponseWriter, r *http.Request) {
		err := updateCache(sheetID, serviceAccountJSON)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Cache updated", "data": cache.Spreadsheet.Data})

	})

	r.Get("/data", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UnixNano() / int64(time.Millisecond)

		if len(cache.Spreadsheet.Data) == 0 || now-cache.Spreadsheet.Timestamp > int64(cacheTTL/time.Millisecond) {
			fmt.Println("Cache expired. Fetching new data...")
			err := updateCache(sheetID, serviceAccountJSON)
			if err != nil {
				http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cache.Spreadsheet.Data)
	})

	fmt.Println("Server running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func updateCache(sheetID, serviceAccountJSON string) error {
	data, err := fetchSheetData(sheetID, serviceAccountJSON)
	if err != nil {
		return err
	}
	cache.Spreadsheet.Data = data
	cache.Spreadsheet.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	return nil
}

func fetchSheetData(sheetID, serviceAccountJSON string) ([]SpreadsheetRow, error) {
	ctx := context.Background()

	conf, err := google.CredentialsFromJSON(ctx, []byte(serviceAccountJSON), sheets.SpreadsheetsReadonlyScope)
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
		return []SpreadsheetRow{}, nil
	}

	var sheetData []SpreadsheetRow
	for _, row := range resp.Values {
		data := SpreadsheetRow{
			Name: fmt.Sprint(row[0]),
		}
		sheetData = append(sheetData, data)
	}

	return sheetData, nil
}
