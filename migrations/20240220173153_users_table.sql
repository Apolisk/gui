-- +goose Up
CREATE TABLE IF NOT EXISTS users (
        name	    TEXT,
        password	TEXT
);

-- +goose Down
DROP TABLE IF EXISTS users;

