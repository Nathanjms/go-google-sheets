package application

import (
	"log/slog"
)

type Config struct {
	BaseURL  string
	HTTPPort int
}

type SpreadsheetRow struct {
	Name string `json:"name"`
}

type Cache struct {
	Spreadsheet struct {
		Data      []SpreadsheetRow `json:"data"`
		Timestamp int64            `json:"timestamp"`
	} `json:"spreadsheet"`
}

type Application struct {
	Config Config
	Logger *slog.Logger
}
