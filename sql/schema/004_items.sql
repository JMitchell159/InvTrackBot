-- +goose Up
CREATE TABLE items(
    name TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    description VARCHAR
);

-- +goose Down
DROP TABLE items;
