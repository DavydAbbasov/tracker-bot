package service

import (
	"context"
	"tracker-bot/internal/models"
)

type LearningRepository interface {
}

type LearningService struct {
	repo LearningRepository
}

func NewLearning(repo LearningRepository) *LearningService {
	return &LearningService{
		repo: repo,
	}
}

func (srv *LearningService) GetLearningStats(ctx context.Context, userID int64) (models.LearningStats, error) {
	return models.LearningStats{
		Language:     "English",
		TotalWords:   463,
		TodayWords:   10,
		LearnedWords: 296,
		NextWordIn:   "23 minutes",
	}, nil
}
