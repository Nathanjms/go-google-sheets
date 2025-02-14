package application

import (
	"log/slog"
	"time"
)

type Config struct {
	BaseURL  string
	HTTPPort int
}

type SpreadsheetRow []interface{}

type SpreadsheetData struct {
	Headers []string         `json:"headers"`
	Data    []SpreadsheetRow `json:"data"`
}

type Cache struct {
	Data struct {
		Spreadsheet struct {
			Data      SpreadsheetData `json:"data"`
			Timestamp int64           `json:"timestamp"`
		} `json:"spreadsheet"`
	}
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
			CacheTTL: 60 * 60 * time.Second,
		},
	}
}
