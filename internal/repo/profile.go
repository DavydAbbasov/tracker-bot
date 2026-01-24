package repo

import (
	"context"
	"errors"
	"tracker-bot/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepo interface {
	Create(ctx context.Context, stats *models.ProfileStats) error
	GetByID(ctx context.Context, id int64) (*models.ProfileStats, error)
	Update(ctx context.Context, id int64, stats *models.ProfileStats) error
	Delete(ctx context.Context, id int64) error
}
type profileRepo struct {
	db *pgxpool.Pool
}

func NewProfileRepo(db *pgxpool.Pool) ProfileRepo {
	return &profileRepo{db: db}
}

func (repo *profileRepo) Create(ctx context.Context, stats *models.ProfileStats) error {
	q := `
		INSERT INTO users (tg_user_id, username, phone_number, email, language, timezone)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := repo.db.Exec(ctx, q, stats.TgUserID, stats.UserName, stats.PhoneNumber, stats.Email, stats.Language, stats.TimeZone)
	if err != nil {
		return err
	}

	return nil
}

func (repo *profileRepo) GetByID(ctx context.Context, id int64) (*models.ProfileStats, error) {
	q := `
	SELECT tg_user_id, username, phone_number, email, language, timezone
	FROM users
	WHERE id = $1
	`
	var profile models.ProfileStats

	err := repo.db.QueryRow(ctx, q, id).Scan(
		&profile.TgUserID,
		&profile.UserName,
		&profile.PhoneNumber,
		&profile.Email,
		&profile.Language,
		&profile.TimeZone,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &profile, nil
}

func (repo *profileRepo) Update(ctx context.Context, id int64, stats *models.ProfileStats) error {
	q := `
		UPDATE users
		SET language = $2, timezone = $3
		WHERE id = $1
	`

	res, err := repo.db.Exec(ctx, q, id, stats.Language, stats.TimeZone)
	if err != nil {
		return err
	}

	rows := res.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (repo *profileRepo) Delete(ctx context.Context, id int64) error {
	q := `
		DELETE FROM users
		WHERE id = $1
	`

	res, err := repo.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	rows := res.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}
