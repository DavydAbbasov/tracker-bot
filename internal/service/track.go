package service

import (
	"context"
	"tracker-bot/internal/models"
	"tracker-bot/internal/repo"
)

type TrackerService interface {
	GetMainStats(ctx context.Context, userID int64) (models.MainStats, error)
}

type trackerService struct {
	repo repo.TrackerRepository
}

func NewTrackerService(repo repo.TrackerRepository) TrackerService {
	return &trackerService{
		repo: repo,
	}
}
func (srv *trackerService) GetMainStats(ctx context.Context, userID int64) (models.MainStats, error) {
	// TODO: заменить на реальные данные из repo
	return models.MainStats{
		CurrentActivityName: "Go",
		TodayTracked:        4*60*60 + 52*60,
		TodaySessions:       4,
		StreakDays:          104,
	}, nil
}
