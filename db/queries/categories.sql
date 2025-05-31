-- name: GetAllCategory :many
SELECT id, 
       shop_id, 
       name,
       parent_id 
FROM categories
LIMIT $1 OFFSET $2;

-- name: GetCategoryByID :one
SELECT id, 
       shop_id, 
       name,
       parent_id 
FROM categories
WHERE id = $1;

-- name: GetCategoriesPaginated :many
SELECT * FROM categories
ORDER BY id
LIMIT $1 OFFSET $2;


-- name: CreateCategory :one
INSERT INTO categories (
    id,
    shop_id,
    name,
    parent_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET 
    shop_id = COALESCE($2, shop_id),
    name = COALESCE($3, name),
    parent_id = COALESCE($4, parent_id)
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;



