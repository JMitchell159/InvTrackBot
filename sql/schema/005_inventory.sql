-- +goose Up
CREATE TABLE inventory(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    quantity INT NOT NULL,
    owner_id UUID NOT NULL REFERENCES players (id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    UNIQUE(owner_id, item_id)
);

-- +goose Down
DROP TABLE inventory;
