package service

import (
	"context"
	"strings"
	"tracker-bot/internal/models"
	"tracker-bot/internal/repo"
)

type TrackerService interface {
	GetMainStats(ctx context.Context, userID int64) (models.MainStats, error)
	CreateActivity(ctx context.Context, userID int64, name, emoji string) (repo.Activity, error)
	ListActivities(ctx context.Context, userID int64) ([]models.TrackActivityItem, error)
	ToggleSelectedActivity(ctx context.Context, userID, activityID int64) error
	DeleteSelectedActivities(ctx context.Context, userID int64) (int64, error)
	ListSelectedActivities(ctx context.Context, userID int64) ([]models.TrackActivityItem, error)
	ListArchivedActivities(ctx context.Context, userID int64) ([]models.TrackActivityItem, error)
	ArchiveSelectedActivities(ctx context.Context, userID int64) (int64, error)
	RestoreArchivedActivity(ctx context.Context, userID, activityID int64) error
	DeleteArchivedForever(ctx context.Context, userID, activityID int64) error
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

func (srv *trackerService) CreateActivity(ctx context.Context, userID int64, name, emoji string) (repo.Activity, error) {
	name = strings.TrimSpace(name)
	emoji = strings.TrimSpace(emoji)
	return srv.repo.Create(ctx, userID, name, emoji)
}

func (srv *trackerService) ListActivities(ctx context.Context, userID int64) ([]models.TrackActivityItem, error) {
	activities, err := srv.repo.ListActive(ctx, userID)
	if err != nil {
		return nil, err
	}

	selectedIDs, err := srv.repo.SelectedListActive(ctx, userID)
	if err != nil {
		return nil, err
	}

	selected := make(map[int64]struct{}, len(selectedIDs))
	for _, id := range selectedIDs {
		selected[id] = struct{}{}
	}

	items := make([]models.TrackActivityItem, 0, len(activities))
	for _, a := range activities {
		_, isSelected := selected[a.ID]
		items = append(items, models.TrackActivityItem{
			ID:       a.ID,
			Name:     a.Name,
			Emoji:    a.Emoji,
			Selected: isSelected,
		})
	}

	return items, nil
}

func (srv *trackerService) ToggleSelectedActivity(ctx context.Context, userID, activityID int64) error {
	return srv.repo.ToggleSelectedActive(ctx, userID, activityID)
}

func (srv *trackerService) DeleteSelectedActivities(ctx context.Context, userID int64) (int64, error) {
	return srv.repo.DeleteSelected(ctx, userID)
}

func (srv *trackerService) ListSelectedActivities(ctx context.Context, userID int64) ([]models.TrackActivityItem, error) {
	items, err := srv.ListActivities(ctx, userID)
	if err != nil {
		return nil, err
	}

	out := make([]models.TrackActivityItem, 0, len(items))
	for _, item := range items {
		if item.Selected {
			out = append(out, item)
		}
	}
	return out, nil
}

func (srv *trackerService) ListArchivedActivities(ctx context.Context, userID int64) ([]models.TrackActivityItem, error) {
	activities, err := srv.repo.ListArchived(ctx, userID)
	if err != nil {
		return nil, err
	}
	items := make([]models.TrackActivityItem, 0, len(activities))
	for _, a := range activities {
		items = append(items, models.TrackActivityItem{
			ID:       a.ID,
			Name:     a.Name,
			Emoji:    a.Emoji,
			Selected: false,
		})
	}
	return items, nil
}

func (srv *trackerService) ArchiveSelectedActivities(ctx context.Context, userID int64) (int64, error) {
	return srv.repo.ArchiveSelected(ctx, userID)
}

func (srv *trackerService) RestoreArchivedActivity(ctx context.Context, userID, activityID int64) error {
	return srv.repo.RestoreArchived(ctx, userID, activityID)
}

func (srv *trackerService) DeleteArchivedForever(ctx context.Context, userID, activityID int64) error {
	return srv.repo.DeleteArchivedForever(ctx, userID, activityID)
}
