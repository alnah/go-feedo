-- +goose Up
CREATE TABLE users (
    id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(100)
);

-- -goose Down
DROP TABLE users;
