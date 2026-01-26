package service

import (
	"context"
	"errors"

	"tracker-bot/internal/models"
	"tracker-bot/internal/repo"
)

type EntryService interface {
	EnsureUser(ctx context.Context, user *models.UserInput) (int64, error)
}

type entryService struct {
	repo repo.EntryRepository
}

func NewEntryService(repo repo.EntryRepository) EntryService {
	return &entryService{
		repo: repo,
	}
}

func (s *entryService) EnsureUser(ctx context.Context, user *models.UserInput) (int64, error) {
	_, err := s.repo.GetByID(ctx, user.TgUserID)
	if err == nil {
		return 0, nil
	}

	if !errors.Is(err, models.ErrUserNotFound) {
		return 0, err
	}

	if user.Language == nil || *user.Language == "" {
		v := "en"
		user.Language = &v
	}
	if user.TimeZone == nil || *user.TimeZone == "" {
		v := "UTC"
		user.TimeZone = &v
	}

	dbID, err := s.repo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrUserExists) {
			return dbID, err
		}
		return 0, err
	}

	return dbID, nil
}
