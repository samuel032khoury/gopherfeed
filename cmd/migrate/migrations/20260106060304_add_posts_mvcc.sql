-- +goose Up
ALTER TABLE posts ADD COLUMN version INT DEFAULT 0 NOT NULL;

-- +goose Down
ALTER TABLE posts DROP COLUMN version;