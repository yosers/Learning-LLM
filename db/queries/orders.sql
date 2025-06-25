-- name: CreateOrder :one
INSERT INTO orders (shop_id, user_id, total, status)
VALUES ($1, $2, $3, $4)
RETURNING id, shop_id, user_id, total, status, created_at;

-- name: GetListOrders :many
SELECT 
  o.id, 
  s.name AS shop_name, 
  u.id AS user_id, 
  o.total, 
  o.status, 
  o.created_at
FROM orders o
JOIN shops s ON o.shop_id = s.id
JOIN users u ON o.user_id = u.id
WHERE ($3::int = 0 OR u.id = $3)
  AND ($4::text = '' OR o.status = $4)
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












