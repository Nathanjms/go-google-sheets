package cache

import (
	"time"

	"github.com/nathanjms/go-google-sheets/internal/application"
)

type Cache struct {
	Spreadsheet struct {
		Data      []application.SpreadsheetRow `json:"data"`
		Timestamp int64                        `json:"timestamp"`
	} `json:"spreadsheet"`
	CacheTTL time.Duration
}

var CacheInstance Cache = Cache{
	CacheTTL: 60 * 60 * time.Second, // Initialize cache TTL
}

var Spreadsheet = CacheInstance.Spreadsheet // Alias for easier access
