package service

import (
	"context"
	"tracker-bot/internal/models"
)

type SubscriptionRepository interface {
}

type SubscriptionService struct {
	repo SubscriptionRepository
}

func NewSubscription(repo SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (srv *SubscriptionService) GetSubscriptionStats(ctx context.Context, userID int64) (models.SubscriptionStats, error) {
	return models.SubscriptionStats{
		ActivePlan: "Free",
		DaysEnd:    23,
	}, nil
}
