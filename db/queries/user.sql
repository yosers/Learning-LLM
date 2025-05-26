-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE shop_id = $1
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (
    id,
    shop_id,
    email,
    phone,
    is_active
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET 
    email = COALESCE($2, email),
    phone = COALESCE($3, phone),
    is_active = COALESCE($4, is_active),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;