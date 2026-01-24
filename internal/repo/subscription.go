package repo

import "github.com/jackc/pgx/v5/pgxpool"

type SubscriptionRepo interface {
}
type subscriptionRepo struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepo(db *pgxpool.Pool) SubscriptionRepo {
	return &subscriptionRepo{db: db}
}
