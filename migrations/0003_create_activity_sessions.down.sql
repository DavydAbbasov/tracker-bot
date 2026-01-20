DROP INDEX IF EXISTS idx_sessions_activity;
DROP INDEX IF EXISTS idx_sessions_user_start_at;
DROP INDEX IF EXISTS uniq_open_session_per_user;
DROP TABLE IF EXISTS activity_sessions;