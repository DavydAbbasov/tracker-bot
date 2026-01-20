UPDATE schema_migrations SET version = 1, dirty = FALSE;
DROP TABLE IF EXISTS activities CASCADE;
DROP INDEX IF EXISTS uq_activities_user_lower_name;
DROP INDEX IF EXISTS idx_activities_user;

DROP INDEX IF EXISTS idx_activities_user;
DROP INDEX IF EXISTS uq_activities_user_lower_name;
DROP TABLE IF EXISTS activities;