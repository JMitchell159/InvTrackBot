-- name: CreateItem :one
INSERT INTO items (name, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: CreateItemWDesc :one
INSERT INTO items (name, created_at, updated_at, description)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: CreateItemWCat :one
INSERT INTO items (name, created_at, updated_at, category)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: CreateItemFull :one
INSERT INTO items (name, created_at, updated_at, category, description)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
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

-- name: GetItemsByCategory :many
SELECT *
FROM items
WHERE category = $1;

-- name: ResetItems :exec
DELETE FROM items;
