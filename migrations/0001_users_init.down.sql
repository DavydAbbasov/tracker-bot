DROP INDEX IF EXISTS idx_users_crreated_at;
DROP INDEX IF EXISTS uniq_phone_number;
DROP INDEX IF EXISTS uniq_users_email;
DROP INDEX IF EXISTS uniq_users_user;
ALTER TABLE IF EXISTS users
    DROP CONSTRAINT IF EXISTS users_tg_user_id_uniq;
DROP TABLE IF EXISTS users;
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_allowed_language;