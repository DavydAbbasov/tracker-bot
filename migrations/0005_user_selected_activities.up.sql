CREATE TABLE IF NOT EXISTS user_selected_activities(
    user_id     BIGINT  REFERENCES users(id) ON DELETE CASCADE,
    activity_id BIGINT  REFERENCES activities(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, activity_id)
);