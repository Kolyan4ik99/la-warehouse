package main

import (
	"io"
	"la-warehouse/config"
	"la-warehouse/internal/repo"
	"la-warehouse/internal/service"
	"la-warehouse/internal/transport"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
)

func main() {
	log := NewLogger()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	go func() {
		<-sig
		log.Info().Msg("application successful exit")
		os.Exit(0)
	}()

	log.Debug().Msg("init config")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Debug().Msg("init postgres connection")
	pg, err := repo.NewPostgres(cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	err = pg.Ping()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Debug().Msg("init service warehouse")
	serv := service.NewWarehouse(log, pg)

	log.Debug().Msg("init handler warehouse")
	wh := transport.NewWarehouse(serv)

	log.Info().Msg("application is started")
	err = wh.ListenAndServe(":" + cfg.Http.Port)
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
