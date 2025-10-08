-- name: CreatePlayer :one
INSERT INTO players (id, created_at, updated_at, name, game_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetPlayer :one
SELECT *
FROM players
WHERE id = $1;

-- name: GetPlayerByName :one
SELECT *
FROM players
WHERE name = $1 AND game_id = $2;

-- name: GetPlayersByGame :many
SELECT *
FROM players
WHERE game_id IN (
    SELECT games.id
    FROM games
    WHERE games.name = $1 AND games.server_id = $2
);
