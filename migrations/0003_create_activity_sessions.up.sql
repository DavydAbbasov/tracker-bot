CREATE TABLE IF NOT EXISTS activity_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,--если пользователя удалят, все его сессии удалятся автоматически.
    activity_id BIGINT NOT NULL REFERENCES activities(id) ON DELETE RESTRICT,--запрещает удалить активность, если по ней есть сессии.
    start_at TIMESTAMPTZ NOT NULL DEFAULT now(),--timestamp with time zone (точка во времени).
    end_at TIMESTAMPTZ NULL,--NULL означает, что сессия ещё открыта.
    planned_min INTEGER NULL,--Плановая длительность в минутах (например, 15/30)
    source TEXT NOT NULL DEFAULT 'prompt',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_planned_min_positive CHECK (planned_min IS NULL OR planned_min > 0),
    CONSTRAINT chk_end_after_start CHECK (end_at IS NULL OR end_at > start_at));
-- Ровно одна "открытая" сессия на пользователя
CREATE UNIQUE INDEX IF NOT EXISTS uniq_open_session_per_user
ON activity_sessions(user_id)--частичный уникальный индекс
WHERE end_at IS NULL;
-- Ускоряем выборки
--Это обычные (неуникальные) индексы
CREATE INDEX IF NOT EXISTS idx_sessions_user_start_at
    ON activity_sessions (user_id, start_at);

CREATE INDEX IF NOT EXISTS idx_sessions_activity
    ON activity_sessions (activity_id);