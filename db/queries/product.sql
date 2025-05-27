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

-- name: ListProducts :many
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
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAllProducts :many
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
WHERE deleted_at IS NULL
ORDER BY created_at DESC; 