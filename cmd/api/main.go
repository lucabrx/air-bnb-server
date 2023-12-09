package main

import (
	"github.com/air-bnb/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

type application struct {
	logger *zerolog.Logger
	wg     sync.WaitGroup
	config config.AppConfig
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	app := application{
		logger: &log.Logger,
		config: cfg,
	}

	err = app.serve()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")

	}
}
