-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE shop_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE shop_id = $1;

-- name: CreateUser :one
INSERT INTO users (
    shop_id,
    email,
    phone,
    is_active
) VALUES (
    $1, $2, $3, $4
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

-- name: FindUserByPhone :one
SELECT * FROM users
WHERE phone = $1 LIMIT 1;

-- name: ListUserRole :many
select rl.* from users us join user_roles ur
on us.id = ur.user_id 
join roles rl on rl.id = ur.role_id 
where us.id = '2'
order by rl.id ASC;