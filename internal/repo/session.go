package repo

import (
	"context"
	"fmt"
	errlocal "tracker-bot/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SessionRepository stores tracked activity sessions.
type SessionRepository interface {
	// CreateRetroSession saves one session ending "now" with duration intervalMin.
	CreateRetroSession(ctx context.Context, userID, activityID int64, intervalMin int, source string) error
}

type sessionRepository struct {
	db *pgxpool.Pool
}

// NewSessionRepository creates session repository backed by pgx pool.
func NewSessionRepository(db *pgxpool.Pool) SessionRepository {
	return &sessionRepository{db: db}
}

// CreateRetroSession writes a backfilled session only for user's active activity.
func (r *sessionRepository) CreateRetroSession(ctx context.Context, userID, activityID int64, intervalMin int, source string) error {
	if userID <= 0 || activityID <= 0 || intervalMin <= 0 {
		return fmt.Errorf("create retro session: invalid input")
	}

	q := `
	INSERT INTO activity_sessions (user_id, activity_id, start_at, end_at, planned_min, source)
	SELECT
		$1,
		$2,
		now() - make_interval(mins => $3),
		now(),
		$3,
		$4
	WHERE EXISTS (
		SELECT 1
		FROM activities
		WHERE id = $2 AND user_id = $1 AND is_archived = FALSE
	);
	`
	// INSERT ... WHERE EXISTS prevents writing sessions for foreign or archived activities.
	tag, err := r.db.Exec(ctx, q, userID, activityID, intervalMin, source)
	if err != nil {
		return fmt.Errorf("create retro session exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errlocal.ErrActivityNotFound
	}
	return nil
}
