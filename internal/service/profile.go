package service

import (
	"context"
	"tracker-bot/internal/models"
)

type ProfileRepository interface {
}

type ProfileService struct {
	repo ProfileRepository
}

func NewProfile(repo ProfileRepository) *ProfileService {
	return &ProfileService{
		repo: repo,
	}
}

func (srv *ProfileService) GetProfileStats(ctx context.Context, userID int64) (models.ProfileStats, error) {
	return models.ProfileStats{
		TgUserID:    userID, // тут надо брать не id из БД
		UserName:    "Alex",
		PhoneNumber: "+142342423",
		Email:       "mock@email.com",
		Language:    "English",
		TimeZone:    "Europe/London",
	}, nil
}
