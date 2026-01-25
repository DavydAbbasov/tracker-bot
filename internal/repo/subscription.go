package repo

import "github.com/jackc/pgx/v5/pgxpool"

type SubscriptionRepository interface {
}
type subscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}
