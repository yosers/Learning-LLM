-- name: CreateOrder :one
INSERT INTO orders (shop_id, user_id, total, status)
VALUES ($1, $2, $3, $4)
RETURNING id, shop_id, user_id, total, status, created_at;

-- name: GetListOrders :many
SELECT id, shop_id, user_id, total, status, created_at
FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetOrderById :one
SELECT id, shop_id, user_id, total, status, created_at
FROM orders
WHERE id = $1;

-- name: UpdateOrder :one
UPDATE orders
SET total = $2, status = $3
WHERE id = $1
RETURNING id, shop_id, user_id, total, status, created_at;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;









