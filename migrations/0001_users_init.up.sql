CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    tg_user_id BIGINT  NOT NULL,
    username CITEXT NULL,
    phone_number TEXT NULL,
    email CITEXT NULL,
    language TEXT NOT NULL DEFAULT 'ru',
    timezone TEXT NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT users_tg_user_id_positive CHECK (tg_user_id > 0),
    CONSTRAINT users_email_format_chk CHECK (email IS NULL OR email ~* '^[^@]+@[^@]+\.[^@]+$'),
    CONSTRAINT users_tg_user_id_uniq UNIQUE (tg_user_id),
    CONSTRAINT users_allowed_language CHECK (language IN('ru','en','de','uk','ar'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_users_username ON users(username) WHERE username IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uniq_users_email ON users(email) WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS uniq_phone_number ON users(phone_number) WHERE phone_number IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_tg ON users(tg_user_id);