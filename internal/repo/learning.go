package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type LearningRepo interface {
}
type learningRepo struct {
	db *pgxpool.Pool
}

func NewLearningRepo(db *pgxpool.Pool) LearningRepo {
	return &learningRepo{db: db}
}
