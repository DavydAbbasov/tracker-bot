package service

import (
	"context"
	"tracker-bot/internal/models"
	"tracker-bot/internal/repo"
)

type LearningService interface {
	GetLearningStats(ctx context.Context, userID int64) (models.LearningStats, error)
}

type learningService struct {
	repo repo.LearningRepository
}

func NewLearningService(repo repo.LearningRepository) LearningService {
	return &learningService{
		repo: repo,
	}
}

func (srv *learningService) GetLearningStats(ctx context.Context, userID int64) (models.LearningStats, error) {
	return models.LearningStats{
		Language:     "English",
		TotalWords:   463,
		TodayWords:   10,
		LearnedWords: 296,
		NextWordIn:   "23 minutes",
	}, nil
}
