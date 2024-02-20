package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"

	"trintech/review/config"
	"trintech/review/pkg/http_server"
	"trintech/review/pkg/postgres_client"
	"trintech/review/pkg/processor"
	"trintech/review/pkg/token_util"
)

var (
	cfgs       *config.Config
	httpServer *http_server.HttpServer
	pgClient   *postgres_client.PostgresClient

	tokenGenerator token_util.JWTAuthenticator

	processors []processor.Processor
	factories  []processor.Factory
)

func loadConfigs() {
	var err error
	cfgs, err = config.LoadConfig("development", os.Getenv("SERVICE"), os.Getenv("ENV"))
	if err != nil {
		log.Fatalln(err.Error())
	}
}

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

func loadTokenGenerator() {
	var err error
	tokenGenerator, err = token_util.NewJWTAuthenticator(cfgs.SymetricKey)
	if err != nil {
		log.Fatalf("unable to create new token generator: %v", err)
	}
}

func loadPostgresClient() {
	pgClient = postgres_client.NewPostgresClient(cfgs.PostgresDB.Address())
}

func loadDefault() {
	loadConfigs()
	loadLogger()
	loadTokenGenerator()
}

func start(ctx context.Context, errChan chan error) {
	for _, f := range factories {
		if err := f.Connect(ctx); err != nil {
			errChan <- err
		}
	}

	for _, p := range processors {
		go func(pr processor.Processor) {
			if err := pr.Start(ctx); err != nil {
				errChan <- err
			}
		}(p)
	}
}

func stop(ctx context.Context) {
	for _, f := range factories {
		if err := f.Close(ctx); err != nil {
			slog.Error("unable to stop factory", "err", err)
		}
	}

	for _, p := range processors {
		if err := p.Stop(ctx); err != nil {
			slog.Error("unable to stop processor", "err", err)
		}
	}
}
