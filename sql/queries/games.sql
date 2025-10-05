-- name: CreateGame :one
INSERT INTO games (id, created_at, updated_at, name, server_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetGame :one
SELECT *
FROM games
WHERE id = $1;

-- name: GetGameByName :one
SELECT *
FROM games
WHERE name = $1 AND server_id = $2;
