package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	errlocal "tracker-bot/internal/errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Activity struct {
	ID         int64
	UserID     int64
	Name       string
	Emoji      string
	IsArchived bool
	CreatedAt  time.Time
}
type TrackerRepository interface {
	Create(ctx context.Context, userID int64, name, emoji string) (Activity, error)
	ListActive(ctx context.Context, userID int64) ([]Activity, error)
	SelectedListActive(ctx context.Context, userID int64) ([]int64, error)
	ToggleSelectedActive(ctx context.Context, userID, activityID int64) error
	DeleteSelected(ctx context.Context, userID int64) (int64, error)
}
type trackRepository struct {
	db *pgxpool.Pool
}

func NewTrackerRepository(db *pgxpool.Pool) TrackerRepository {
	return &trackRepository{db: db}
}

func (r *trackRepository) Create(ctx context.Context, userID int64, name, emoji string) (Activity, error) {
	if userID <= 0 {
		return Activity{}, fmt.Errorf("create activity: invalid userID")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return Activity{}, fmt.Errorf("create activity: empty name")
	}

	emoji = strings.TrimSpace(emoji)

	q := `
	INSERT INTO activities (user_id, name, emoji)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, name, emoji, is_archived, created_at;
	`

	var a Activity
	err := r.db.QueryRow(ctx, q, userID, name, emoji).Scan(
		&a.ID,
		&a.UserID,
		&a.Name,
		&a.Emoji,
		&a.IsArchived,
		&a.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23505 = unique_violation
			if pgErr.Code == "23505" {
				return Activity{}, errlocal.ErrActivityExists
			}
		}
		return Activity{}, fmt.Errorf("create activity: %w", err)
	}
	return a, nil
}

func (r *trackRepository) ListActive(ctx context.Context, userID int64) ([]Activity, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("list active: invalid userID")
	}

	q := `
	SELECT id, user_id, name, emoji, is_archived, created_at
	FROM activities
	WHERE user_id = $1 AND is_archived = false
	ORDER BY lower(name), id;
	`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list active query: %w", err)
	}
	defer rows.Close()

	out := make([]Activity, 0, 16)
	for rows.Next() {
		var a Activity
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Name, &a.Emoji, &a.IsArchived, &a.CreatedAt,
		); err != nil {
		}
		out = append(out, a)
		return nil, fmt.Errorf("list active scan: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list active rows: %w", err)
	}
	return out, nil
}

func (r *trackRepository) SelectedListActive(ctx context.Context, userID int64) ([]int64, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("selected list: invalid userID")
	}

	q := `
	SELECT activity_id
	FROM user_selected_activities
	WHERE user_id = $1
	ORDER BY activity_id;
	`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("selected list query: %w", err)
	}
	defer rows.Close()

	ids := make([]int64, 0, 16)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("selected list scan: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("selected list rows: %w", err)
	}
	return ids, nil
}

// to do transaction
func (r *trackRepository) ToggleSelectedActive(ctx context.Context, userID, activityID int64) error {
	if userID <= 0 || activityID <= 0 {
		return fmt.Errorf("toggle selected: invalid ids")
	}

	// В транзакции, чтобы не было гонок между проверкой владения и insert/delete
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("toggle selected begin: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1) Проверка владения активностью
	ownQ := `
	SELECT EXISTS(
		SELECT 1
		FROM activities
		WHERE id = $1 AND user_id = $2 AND is_archived = false
	);`
	var owned bool
	if err := tx.QueryRow(ctx, ownQ, activityID, userID).Scan(&owned); err != nil {
		return fmt.Errorf("toggle selected ownership: %w", err)
	}
	if !owned {
		return errlocal.ErrActivityNotFound
	}

	// 2) Пробуем удалить (если было выбрано)
	delQ := `
	DELETE FROM user_selected_activities
	WHERE user_id = $1 AND activity_id = $2;`
	tag, err := tx.Exec(ctx, delQ, userID, activityID)
	if err != nil {
		return fmt.Errorf("toggle selected delete: %w", err)
	}
	if tag.RowsAffected() == 1 {
		return tx.Commit(ctx)
	}

	// 3) Не было — вставляем
	insQ := `
	INSERT INTO user_selected_activities(user_id, activity_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING;`
	if _, err := tx.Exec(ctx, insQ, userID, activityID); err != nil {
		return fmt.Errorf("toggle selected insert: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *trackRepository) DeleteSelected(ctx context.Context, userID int64) (int64, error) {
	if userID <= 0 {
		return 0, fmt.Errorf("delete selected: invalid userID")
	}

	q := `
	DELETE FROM activities a
	USING user_selected_activities s
	WHERE a.id = s.activity_id
	  AND a.user_id = $1
	  AND s.user_id = $1;
	`

	tag, err := r.db.Exec(ctx, q, userID)
	if err != nil {
		return 0, fmt.Errorf("delete selected exec: %w", err)
	}
	return tag.RowsAffected(), nil
}
