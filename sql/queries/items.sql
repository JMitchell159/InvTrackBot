-- name: CreateItem :one
INSERT INTO items (name, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: UpdateDesc :exec
UPDATE items
SET description = $1, updated_at = $2
WHERE name = $3;

-- name: UpdateCat :exec
UPDATE items
SET category = $1, updated_at = $2
WHERE name = $3;

-- name: GetItem :one
SELECT *
FROM items
WHERE name = $1;

-- name: ResetItems :exec
DELETE FROM items;
