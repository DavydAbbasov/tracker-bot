package repo

import "github.com/jackc/pgx/v5/pgxpool"

type EntryRepository interface {
}
type entryRepository struct {
	db *pgxpool.Pool
}

func NewEntryRepository(db *pgxpool.Pool) EntryRepository {
	return &entryRepository{db: db}
}
