package main

import (
	"io"
	"la-warehouse/config"
	"la-warehouse/internal/repo"
	"la-warehouse/internal/service"
	"la-warehouse/internal/transport"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func main() {
	log := NewLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	pg, err := repo.NewPostgres(cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	err = pg.Ping()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	serv := service.NewWarehouse(log, pg)

	wh := transport.NewWarehouse(serv)
	err = wh.ListenAndServe(":8080")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func NewLogger() *zerolog.Logger {
	mw := io.MultiWriter(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006.01.02 15:04:05.000",
	})

	zerolog.TimeFieldFormat = time.RFC3339Nano

	logger := zerolog.New(mw).With().Timestamp().Logger()

	logger = logger.Level(zerolog.DebugLevel)

	logger.Info().Msg("logger loaded")

	return &logger
}
