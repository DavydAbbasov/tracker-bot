package main

import (
	"context"
	"os/signal"
	"syscall"
	"tracker-bot/internal/application"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	app := application.NewApplication()
	if err := app.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("All systems closed with errors!")
	}

}
