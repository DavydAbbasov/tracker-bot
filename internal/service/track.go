package service

import (
	"context"
	"tracker-bot/internal/models"
)

type TrackerRepository interface {
	EnsureIDByTelegram(ctx context.Context, tgID int64, username string) (int64, error)
	InsertUser(ctx context.Context, tgID int64, username *string) error
	GetUserByTelegramID(ctx context.Context, tgID int64) (*models.User, error)
	UpdateUsername(ctx context.Context, tgID int64, uname string) (int64, error)
	UpdateLanguage(ctx context.Context, tgID int64, lang string) (int64, error)
}

type TrackerService struct {
	repo TrackerRepository
}

func NewTracker(repo TrackerRepository) *TrackerService {
	return &TrackerService{
		repo: repo,
	}
}
func (srv *TrackerService) GetMainStats(ctx context.Context, userID int64) (models.MainStats, error) {
	// TODO: заменить на реальные данные из repo
	return models.MainStats{
		CurrentActivityName: "Go",
		TodayTracked:        4*60*60 + 52*60,
		TodaySessions:       4,
		StreakDays:          104,
	}, nil
}
