package application

import (
	"context"
	"fmt"
	"tracker-bot/internal/config"
	"tracker-bot/internal/dispatcher"
	router "tracker-bot/internal/handlers"
	"tracker-bot/internal/repo"
	"tracker-bot/internal/scheduler"
	"tracker-bot/internal/service"
	"tracker-bot/internal/utils/pgclient"
	tgclient "tracker-bot/internal/utils/tgcient"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Application struct {
	cfg        *config.Config
	db         *pgclient.PostgreDB
	repo       *repo.TrackerRepository
	service    *service.EntryService
	bot        *tgbotapi.BotAPI
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

	entryRepo := repo.NewEntryRepository(app.db.Pool())
	profileRepo := repo.NewProfileRepository(app.db.Pool())
	trackRepo := repo.NewTrackerRepository(app.db.Pool())
	learningRepo := repo.NewLearningRepository(app.db.Pool())
	subscriptionRepo := repo.NewSubscriptionRepository(app.db.Pool())
	timerRepo := repo.NewTimerRepository(app.db.Pool())
	sessionRepo := repo.NewSessionRepository(app.db.Pool())

	entrysvc := service.NewEntryService(entryRepo)
	provilesvc := service.NewProfileService(profileRepo)
	tracksvc := service.NewTrackerService(trackRepo)
	timersvc := service.NewTimerService(timerRepo, sessionRepo)
	learningsvc := service.NewLearningService(learningRepo)
	subscriptionsvc := service.NewSubscriptionService(subscriptionRepo)

	module := router.New(app.bot, entrysvc, provilesvc, tracksvc, timersvc, learningsvc, subscriptionsvc, app.cfg.TestTimerMinutes)
	app.dispatcher = dispatcher.New(app.bot, ctx, entrysvc, module, module, module, module, module)
	timerScheduler := scheduler.NewTimerScheduler(ctx, timersvc, module)
	timerScheduler.Run()

	app.dispatcher.Run()
	fmt.Println("dispatcher running")

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
