-- +goose Up
CREATE TABLE games(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    server_id TEXT NOT NULL REFERENCES servers (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE games;
