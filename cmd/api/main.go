package main

import (
	"context"
	"database/sql"
	"github.com/air-bnb/config"
	"github.com/air-bnb/internal/data"
	"github.com/air-bnb/internal/mailer"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
	"time"
)

type application struct {
	logger *zerolog.Logger
	wg     sync.WaitGroup
	config config.AppConfig
	models data.Models
	mailer mailer.Mailer
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	db, err := dbConnection(cfg.DBUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()
	log.Logger.Info().Msg("Connected to database")

	app := application{
		logger: &log.Logger,
		config: cfg,
		models: data.NewModels(db),
		mailer: mailer.NewMailer(cfg.ResendApiKey),
	}

	err = app.serve()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")

	}
}

func dbConnection(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(15 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
