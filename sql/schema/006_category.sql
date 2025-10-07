-- +goose Up
ALTER TABLE items
ADD category TEXT;

-- +goose Down
ALTER TABLE items
DROP COLUMN category;
