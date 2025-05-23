package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/nathanjms/go-google-sheets/internal/application"
	"github.com/nathanjms/go-google-sheets/internal/env"
	"github.com/nathanjms/go-google-sheets/internal/version"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)

		// Spin up a basic http server for the 500 error:
		httpErr := http.ListenAndServe(fmt.Sprintf(":%d", env.GetInt("PORT", 3000)), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
		}))
		if httpErr != nil {
			log.Fatal(httpErr)
			os.Exit(1)
		}

	}
}

func run(logger *slog.Logger) error {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file, proceeding with system environment variables")
	}

	fmt.Printf("version: %s\n", version.Get())

	// --- APP ---
	app := &application.Application{
		Config: application.Config{
			BaseURL:              env.GetString("BASE_URL", "http://localhost"),
			HTTPPort:             env.GetInt("PORT", 3000),
			SpreadsheetId:        env.GetString("GOOGLE_SHEET_ID", ""),
			GoogleServiceAccount: env.GetString("GOOGLE_SERVICE_ACCOUNT", ""),
		},
		Logger: logger,
		Cache: application.Cache{ // Initialize the cache here
			CacheTTL: 60 * 60 * time.Second,
		},
	}

	return serveHttp(app)
}
