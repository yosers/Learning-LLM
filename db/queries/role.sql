-- name: ListRole :many
SELECT * FROM roles
where is_active = true;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1 and is_active = true LIMIT 1;


-- name: DeleteRoleById :exec
UPDATE roles
SET is_active = false, updated_at = now()
WHERE id = $1;

-- name: UpdateRoleById :exec
UPDATE roles
SET name = $2, is_active = $3, updated_at = now()
WHERE id = $1;

-- name: GetRoleByName :one
SELECT * FROM roles
WHERE name = $1 and is_active = true LIMIT 1;

-- name: CreateRole :one
INSERT INTO roles (
    name,
    is_active
) VALUES (
    $1, $2
)   
RETURNING *;