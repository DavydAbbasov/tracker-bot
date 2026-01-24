package application

import (
	"context"
	"fmt"
	"tracker-bot/internal/config"
	"tracker-bot/internal/dispatcher"
	router "tracker-bot/internal/handlers"
	"tracker-bot/internal/repo"
	"tracker-bot/internal/service"
	"tracker-bot/internal/utils/pgclient"
	tgclient "tracker-bot/internal/utils/tgcient"
)

type Application struct {
	cfg        *config.Config
	db         *pgclient.PostgreDB
	repo       *repo.TrackerRepo
	service    *service.EntryService
	bot        *tgclient.Client
	dispatcher *dispatcher.Dispatcher
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
	app.db = db

	bot, err := tgclient.New(app.cfg.Telegram.TelegramToken)
	if err != nil {
		return fmt.Errorf("init telegram bot: %w", err)
	}
	app.bot = bot

	profileRepo := repo.NewProfileRepo(app.db.Pool())
	trackRepo := repo.NewTrackRepo(app.db.Pool())
	learningRepo := repo.NewLearningRepo(app.db.Pool())
	subscriptionRepo := repo.NewSubscriptionRepo(app.db.Pool())

	provilesvc := service.NewProfileService(profileRepo)
	tracksvc := service.NewTracker(trackRepo)
	learningsvc := service.NewLearning(learningRepo)
	subscriptionsvc := service.NewSubscription(subscriptionRepo)

	module = router.New(app.bot, provilesvc, tracksvc, learningsvc, subscriptionsvc)

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
