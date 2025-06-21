-- name: CreateOrder :one
INSERT INTO orders (shop_id, user_id, total, status)
VALUES ($1, $2, $3, $4)
RETURNING id, shop_id, user_id, total, status, created_at;

-- name: GetListOrders :many
SELECT o.id, s.name as shop_name, u.id as user_id, o.total, o.status, o.created_at
FROM orders o inner join shops s on o.shop_id = s.id
inner join users u on o.user_id = u.id
ORDER BY o.created_at DESC
LIMIT $1 OFFSET $2;


-- name: GetOrderById :one
SELECT id, shop_id, user_id, total, status, created_at
FROM orders
WHERE id = $1;

-- name: UpdateOrder :one
UPDATE orders
SET total = COALESCE($2, total), status = COALESCE($3, status)  
WHERE id = $1
RETURNING id, shop_id, user_id, total, status, created_at;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: GetCountOrder :one
SELECT COUNT(*) 
FROM orders;












