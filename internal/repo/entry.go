package repo

import (
	"context"
	"errors"
	"tracker-bot/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EntryRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Create(ctx context.Context, stats *models.UserInput) (int64, error)
}
type entryRepository struct {
	db *pgxpool.Pool
}

func NewEntryRepository(db *pgxpool.Pool) EntryRepository {
	return &entryRepository{db: db}
}

func (repo *entryRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	q := `
	SELECT tg_user_id, username, phone_number, email, language, timezone
	FROM users
	WHERE tg_user_id = $1
	`
	var user models.User

	err := repo.db.QueryRow(ctx, q, id).Scan(
		&user.TgUserID,
		&user.UserName,
		&user.PhoneNumber,
		&user.Email,
		&user.Language,
		&user.TimeZone,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *entryRepository) Create(ctx context.Context, user *models.UserInput) (int64, error) {
	q := `
		INSERT INTO users (tg_user_id, username, phone_number, email, language, timezone)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id int64
	err := r.db.QueryRow(ctx, q,
		user.TgUserID,
		user.UserName,
		user.PhoneNumber,
		user.Email,
		user.Language,
		user.TimeZone,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
