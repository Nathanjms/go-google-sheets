package application

import (
	"log/slog"
	"time"
)

type Config struct {
	BaseURL              string
	HTTPPort             int
	SpreadsheetId        string
	GoogleServiceAccount string
}

type SpreadsheetRow []interface{}

type SpreadsheetData struct {
	Headers  []string         `json:"headers"`
	Contents []SpreadsheetRow `json:"contents"`
}

type SpreadsheetCache struct {
	Data      SpreadsheetData `json:"data"`
	Timestamp int64           `json:"timestamp"`
}

type Spreadsheets map[string]map[string]SpreadsheetCache

type Cache struct {
	Data struct {
		Spreadsheets Spreadsheets `json:"spreadsheets"`
	} `json:"data"`
	CacheTTL time.Duration
}

type Application struct {
	Config Config
	Logger *slog.Logger
	Cache  Cache // Add the Cache field
}

// NewApplication is a constructor for the application.
func NewApplication(cfg Config, logger *slog.Logger) *Application {
	return &Application{
		Config: cfg,
		Logger: logger,
		Cache: Cache{ // Initialize the cache here
			Data: struct {
				Spreadsheets Spreadsheets `json:"spreadsheets"`
			}{Spreadsheets: make(Spreadsheets)},
			CacheTTL: 60 * 60 * time.Second,
		},
	}
}
