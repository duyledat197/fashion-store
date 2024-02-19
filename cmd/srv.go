package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func loadLogger() {
	var logger *slog.Logger
	switch os.Getenv("ENV") {
	case "dev":
		logger = slog.New(tint.NewHandler(os.Stdout, nil))
	case "stg", "prd":
		output := os.Getenv("FILE_LOG_OUTPUT")
		f, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("unable to open log file output: %w", err))
		}
		logger = slog.New(slog.NewJSONHandler(f, nil))
	}
	slog.SetDefault(logger)
}
