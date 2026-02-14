package main

import (
	"context"
	"os/signal"
	"syscall"
	"tracker-bot/internal/application"
	"tracker-bot/internal/config"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}

	app := application.NewApplication(cfg)
	if err := app.Build(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to build application")
	}

	if err := app.Run(); err != nil {
		log.Fatal().Err(err).Msg("All systems closed with errors!")
	}

}
