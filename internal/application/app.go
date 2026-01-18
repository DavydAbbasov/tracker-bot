package application

import (
	"context"
	"fmt"
	"tracker-bot/internal/config"
	"tracker-bot/internal/repo"
	"tracker-bot/internal/service"
	"tracker-bot/internal/utils/pgclient"
)

type Application struct {
	cfg     *config.Config
	db      *pgclient.PostgreDB
	repo    *repo.Repo
	service *service.EntryService
}

func NewApplication() *Application {
	return &Application{}
}
func (app *Application) Start(ctx context.Context) error {
	if err := app.initConfig(); err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	dsn := app.cfg.PostgresDSN()
	db, err := pgclient.NewPgProvider(ctx, dsn)
	if err != nil {
		return fmt.Errorf("init pg client: %w", err)
	}
	_ = db
	return nil
}
func (a *Application) initConfig() error {
	cfg, err := config.ParseConfig()
	if err != nil {
		return err
	}

	a.cfg = cfg
	return nil
}
