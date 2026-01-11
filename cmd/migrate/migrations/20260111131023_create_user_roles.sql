-- +goose Up
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, level, description) VALUES
('user', 1, 'Standard user with basic access'),
('moderator', 4, 'Moderator with elevated privileges'),
('admin', 10, 'Administrator with full access');

-- +goose Down
DROP TABLE IF EXISTS roles;
