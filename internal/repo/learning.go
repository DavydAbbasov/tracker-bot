package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type LearningRepository interface {
}
type learningRepository struct {
	db *pgxpool.Pool
}

func NewLearningRepository(db *pgxpool.Pool) LearningRepository {
	return &learningRepository{db: db}
}
