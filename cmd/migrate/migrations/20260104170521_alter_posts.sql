-- +goose Up
ALTER TABLE posts
    ADD CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE;

-- +goose Down
ALTER TABLE posts
    DROP CONSTRAINT fk_user;