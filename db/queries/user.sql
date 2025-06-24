-- name: GetUser :one
SELECT us.*, s."name" FROM users us join shops s on us.shop_id = s.id 
WHERE us.id = $1 and us.is_active = true LIMIT 1;

-- name: ListUsers :many
SELECT us.*, s."name" as shopName FROM users us join shops s on us.shop_id = s.id 
WHERE  us.is_active = true
ORDER BY us.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE is_active = true;
-- WHERE shop_id = $1;

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
WHERE phone = $1 AND is_active = true LIMIT 1;

-- name: ListUserRole :many
select rl.* from users us join user_roles ur
on us.id = ur.user_id 
join roles rl on rl.id = ur.role_id 
where us.id = $1 AND us.is_active = true
order by rl.id ASC;

-- name: FindUserByPhoneAndCode :one
SELECT * FROM users
WHERE phone = $1 and code_area = $2  AND is_active = true LIMIT 1;

-- name: DeleteUserById :exec
UPDATE users    
SET is_active = false, updated_at = now()
WHERE id = $1;