-- name: CreateUserProfile :one
INSERT INTO user_profiles (
    user_id,
    phone,
    first_name,
    last_name,
    address,
    city,
    country,
    postal_code
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetUserProfile :one
SELECT * FROM user_profiles
WHERE user_id = $1 LIMIT 1;

-- name: UpdateUserProfile :exec
UPDATE user_profiles
SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    address = COALESCE($4, address),
    city = COALESCE($5, city),
    country = COALESCE($6, country),
    postal_code = COALESCE($7, postal_code),
    phone = COALESCE($8, phone)
WHERE user_id = $1; 