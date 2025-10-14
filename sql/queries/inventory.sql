-- name: AddLineItem :one
WITH inserted_line_item AS (
    INSERT INTO inventory (id, created_at, updated_at, quantity, owner_id, item_name)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
    RETURNING *
)
SELECT inserted_line_item.*, players.name AS owner_name
FROM inserted_line_item
INNER JOIN players
ON inserted_line_item.owner_id = players.id
INNER JOIN items
ON inserted_line_item.item_name = items.name;

-- name: UpdateLineItemWName :exec
UPDATE inventory
SET quantity = quantity + $1, updated_at = $2
WHERE owner_id IN (
    SELECT players.id
    FROM players
    WHERE players.name = $3 AND players.game_id IN (
        SELECT games.id
        FROM games
        WHERE games.name = $4 AND games.server_id = $5
    )
) AND item_name = $6;

-- name: UpdateLineItemWID :exec
UPDATE inventory
SET quantity = quantity + $1, updated_at = $2
WHERE owner_id = $3 AND item_name = $4;

-- name: GetItemsByOwnerName :many
SELECT items.*, inventory.quantity AS quantity
FROM inventory
INNER JOIN players
ON inventory.owner_id = players.id
INNER JOIN items
ON inventory.item_name = items.name
WHERE inventory.owner_id IN (
    SELECT players.id
    FROM players
    WHERE players.name = $1 AND players.game_id IN (
        SELECT games.id
        FROM games
        WHERE games.name = $2 AND games.server_id = $3
    )
);

-- name: GetItemsByOwner :many
SELECT items.*, inventory.quantity AS quantity
FROM inventory
INNER JOIN players
ON inventory.owner_id = players.id
INNER JOIN items
ON inventory.item_name = items.name
WHERE inventory.owner_id = $1;

-- name: GetLineItemByItemAndOwner :one
SELECT inventory.*
FROM inventory
INNER JOIN players
ON inventory.owner_id = players.id
INNER JOIN items
ON inventory.item_name = items.name
WHERE inventory.owner_id = $1 AND inventory.item_name = $2;

-- name: GetLineItemByItemAndOwnerName :one
SELECT inventory.*
FROM inventory
INNER JOIN players
ON inventory.owner_id = players.id
INNER JOIN items
ON inventory.item_name = items.name
WHERE inventory.owner_id IN (
    SELECT players.id
    FROM players
    WHERE players.name = $1 AND players.game_id IN (
        SELECT games.id
        FROM games
        WHERE games.name = $2 AND games.server_id = $3
    )
) AND inventory.item_name = $4;
