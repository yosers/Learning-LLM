-- name: CreateOrderItems :one
INSERT INTO order_items (id, order_id, product_id, quantity, unit_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, order_id, product_id, quantity, unit_price;


-- name: GetOrderItemsByID :many
select o.id as order_id, pr.name, pr.description, oi.quantity, oi.unit_price from orders o 
join order_items oi on o.id = oi.order_id
join products pr on pr.id = oi.product_id 
where o.id = $1;





