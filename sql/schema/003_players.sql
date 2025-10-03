-- +goose Up
CREATE TABLE players(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    game_id UUID NOT NULL REFERENCES games (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE players;
