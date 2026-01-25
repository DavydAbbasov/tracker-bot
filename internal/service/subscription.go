package service

import (
	"context"
	"tracker-bot/internal/models"
	"tracker-bot/internal/repo"
)

type SubscriptionService interface {
	GetSubscriptionStats(ctx context.Context, userID int64) (models.SubscriptionStats, error)
}

type subscriptionService struct {
	repo repo.SubscriptionRepository
}

func NewSubscriptionService(repo repo.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		repo: repo,
	}
}

func (srv *subscriptionService) GetSubscriptionStats(ctx context.Context, userID int64) (models.SubscriptionStats, error) {
	return models.SubscriptionStats{
		ActivePlan: "Free",
		DaysEnd:    23,
	}, nil
}
