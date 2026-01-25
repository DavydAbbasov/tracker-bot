package service

import (
	"context"
	"tracker-bot/internal/models"
	"tracker-bot/internal/repo"
)

type ProfileService interface {
	GetProfileStats(ctx context.Context, userID int64) (*models.ProfileStats, error)
	ChangeLanguage(ctx context.Context, userID int64, language string) error
	ChangeTimeZone(ctx context.Context, userID int64, timezone string) error
}

type profileService struct {
	repo repo.ProfileRepository
}

func NewProfileService(repo repo.ProfileRepository) ProfileService {
	return &profileService{
		repo: repo,
	}
}

func (srv *profileService) GetProfileStats(ctx context.Context, userID int64) (*models.ProfileStats, error) {
	profile, err := srv.repo.GetByID(ctx, userID)
	if err != nil {
		return &models.ProfileStats{}, err
	}

	return profile, nil
}

func (srv *profileService) ChangeLanguage(ctx context.Context, userID int64, language string) error {
	profile, err := srv.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	profile.Language = language

	err = srv.repo.Update(ctx, userID, profile)
	if err != nil {
		return err
	}

	return nil
}

func (srv *profileService) ChangeTimeZone(ctx context.Context, userID int64, timezone string) error {
	profile, err := srv.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	profile.TimeZone = timezone

	err = srv.repo.Update(ctx, userID, profile)
	if err != nil {
		return err
	}

	return nil
}
