-- +goose Up
CREATE TABLE IF NOT EXISTS user_invitations (
    token bytea PRIMARY KEY NOT NULL,
    user_id bigint REFERENCES users(id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL
);

ALTER TABLE users
ADD COLUMN IF NOT EXISTS is_active boolean DEFAULT FALSE NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS user_invitations;
ALTER TABLE users
DROP COLUMN IF EXISTS is_active;
