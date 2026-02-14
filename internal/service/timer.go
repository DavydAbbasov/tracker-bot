package service

import (
	"context"
	"fmt"
	"time"
	"tracker-bot/internal/models"
	"tracker-bot/internal/repo"
)

// TimerService contains timer-related use-cases.
type TimerService interface {
	Activate(ctx context.Context, userID int64, intervalMin int) error
	Stop(ctx context.Context, userID int64) error
	ListDueUsers(ctx context.Context, now time.Time, limit int) ([]models.TimerDueUser, error)
	MarkPromptSent(ctx context.Context, userID int64, intervalMin int, now time.Time) error
	RecordPromptAnswer(ctx context.Context, userID, activityID int64) error
	RecordPromptAnswerWithInterval(ctx context.Context, userID, activityID int64, intervalMin int) error
}

type timerService struct {
	timerRepo   repo.TimerRepository
	sessionRepo repo.SessionRepository
}

// NewTimerService creates timer service.
func NewTimerService(timerRepo repo.TimerRepository, sessionRepo repo.SessionRepository) TimerService {
	return &timerService{
		timerRepo:   timerRepo,
		sessionRepo: sessionRepo,
	}
}

// Activate enables timer and schedules next prompt.
func (s *timerService) Activate(ctx context.Context, userID int64, intervalMin int) error {
	if userID <= 0 {
		return fmt.Errorf("activate timer: invalid userID")
	}
	if intervalMin <= 0 {
		return fmt.Errorf("activate timer: invalid interval")
	}
	nextPingAt := time.Now().UTC().Add(time.Duration(intervalMin) * time.Minute)
	return s.timerRepo.UpsertInterval(ctx, userID, intervalMin, nextPingAt)
}

// Stop disables timer for user.
func (s *timerService) Stop(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return fmt.Errorf("stop timer: invalid userID")
	}
	return s.timerRepo.Disable(ctx, userID)
}

// ListDueUsers returns users that should receive prompt now.
func (s *timerService) ListDueUsers(ctx context.Context, now time.Time, limit int) ([]models.TimerDueUser, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.timerRepo.ListDueUsers(ctx, now.UTC(), limit)
}

// MarkPromptSent moves next prompt time forward by interval.
func (s *timerService) MarkPromptSent(ctx context.Context, userID int64, intervalMin int, now time.Time) error {
	nextPingAt := now.UTC().Add(time.Duration(intervalMin) * time.Minute)
	return s.timerRepo.SetNextPing(ctx, userID, nextPingAt)
}

// RecordPromptAnswer stores prompt answer using current timer interval from settings.
func (s *timerService) RecordPromptAnswer(ctx context.Context, userID, activityID int64) error {
	intervalMin, err := s.timerRepo.GetInterval(ctx, userID)
	if err != nil {
		return fmt.Errorf("get interval: %w", err)
	}
	return s.sessionRepo.CreateRetroSession(ctx, userID, activityID, intervalMin, "prompt")
}

// RecordPromptAnswerWithInterval stores prompt answer for explicit interval.
func (s *timerService) RecordPromptAnswerWithInterval(ctx context.Context, userID, activityID int64, intervalMin int) error {
	if intervalMin <= 0 {
		return fmt.Errorf("invalid interval")
	}
	return s.sessionRepo.CreateRetroSession(ctx, userID, activityID, intervalMin, "prompt")
}
