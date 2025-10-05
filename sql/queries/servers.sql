-- name: CreateServer :one
INSERT INTO servers (id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetServer :one
SELECT *
FROM servers
WHERE id = $1;
