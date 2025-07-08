-- name: ListShops :many
SELECT s.* FROM shops s 
WHERE  s.is_active = true
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountShops :one
SELECT COUNT(*) FROM shops
WHERE is_active = true;

-- name: GetShopsById :one
SELECT s.* FROM shops s 
WHERE s.id = $1 and s.is_active = true LIMIT 1;

-- name: UpdateShops :one
UPDATE shops
SET 
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    logo_url = COALESCE($4, logo_url),
    website_url = COALESCE($5, website_url),
    email = COALESCE($6, email),
    whatsapp_phone = COALESCE($7, whatsapp_phone),
    address = COALESCE($8, address),
    city = COALESCE($9, city),
    state = COALESCE($10, state),
    zip_code = COALESCE($11, zip_code),
    country = COALESCE($12, country),
    latitude = COALESCE($13, latitude),
    longitude = COALESCE($14, longitude),
    is_active = COALESCE($15, is_active)
WHERE id = $1
RETURNING *;

-- name: DeleteShopsById :exec
UPDATE shops    
SET is_active = false, updated_at = now()
WHERE id = $1;

-- name: GetShopsByNameOrWhatshapp :one
SELECT s.* FROM shops s 
WHERE s.is_active = true 
  AND (
    s.name = $1 
    OR s.whatsapp_phone = $2
) LIMIT 1;

-- name: CreateShops :one
INSERT INTO shops (
    name,
    description,
    logo_url,
    website_url,
    email,
    whatsapp_phone,
    address,
    city,
    state,
    zip_code,
    country,
    latitude,
    longitude,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING *;

-- name: GetAllShops :many
SELECT s.* FROM shops s 
WHERE  s.is_active = true
ORDER BY s.created_at DESC;