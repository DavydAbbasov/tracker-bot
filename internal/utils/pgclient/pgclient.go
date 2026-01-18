package pgclient

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgreDB struct {
	db *pgxpool.Pool
}

func NewPgProvider(ctx context.Context, dsn string) (*PostgreDB, error) {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error get database driver %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error connection to database %w", err)
	}

	log.Info().Msg("database Repos connection is success")

	r := &PostgreDB{
		db: db,
	}
	return r, nil
}
func (c *PostgreDB) Pool() *pgxpool.Pool {
	return c.db
}

func (c *PostgreDB) Close() {
	if c.db != nil {
		c.db.Close()
	}
}
