--1:1
-- хранит персональные настройки напоминаний (пингов)
-- для каждого пользователя и точку времени следующего пинга.
CREATE TABLE IF NOT EXISTS user_timer_settings (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,-- CASCSDE дочерние записи не живут без родителя
    interval_min INTEGER NOT NULL DEFAULT 15,
    next_ping_at TIMESTAMPTZ NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,--Главный тумблер: включены ли пинги для пользователя.
    timezone      TEXT        NOT NULL DEFAULT 'Europe/Berlin',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_interval_min_range CHECK (interval_min > 0 AND interval_min <= 360)
);
--Частичный индекс: ускоряет выборку только включённых пользователей
--по времени next_ping_at — бот быстро найдёт, «кому слать сейчас».
CREATE INDEX IF NOT EXISTS idx_timer_due
    ON user_timer_settings(next_ping_at)
    WHERE enabled = TRUE;