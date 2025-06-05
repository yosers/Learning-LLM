-- name: GetProductByID :one
SELECT id, 
       name, 
       description, 
       price, 
       stock, 
       category_id,
       created_at, 
       updated_at, 
       deleted_at
FROM products
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateProduct :one
INSERT INTO products (id, name, description, price, stock, category_id, shop_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, description, price, stock, category_id, shop_id, created_at, updated_at, deleted_at;

-- name: UpdateProduct :one
UPDATE products
SET name = COALESCE($2, name), description = COALESCE($3, description), 
       price = COALESCE($4, price), stock = COALESCE($5, stock), category_id = COALESCE($6, category_id),
       shop_id = COALESCE($7, shop_id)
WHERE id = $1
RETURNING *;

-- name: ListProducts :many
SELECT p.id, 
       p.name, 
       p.description, 
       p.price, 
       p.stock, 
       c.name as category_id,
       s.name as shop_id,
       p.created_at, 
       p.updated_at, 
       p.deleted_at
FROM products p inner join categories c on p.category_id = c.id
inner join shops s on p.shop_id = s.id
WHERE p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAllProducts :many
SELECT p.id, 
       p.name, 
       p.description, 
       p.price, 
       p.stock, 
       c.name as category_id,
       s.name as shop_id,
       p.created_at, 
       p.updated_at, 
       p.deleted_at
FROM products p inner join categories c on p.category_id = c.id
inner join shops s on p.shop_id = s.id
WHERE p.deleted_at IS NULL
ORDER BY p.created_at DESC; 

-- name: DeleteProductByID :exec
UPDATE products
SET deleted_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetCountProduct :one
SELECT COUNT(*) 
FROM products 
WHERE deleted_at IS NULL;