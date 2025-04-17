-- name: GetClient :one
SELECT * FROM my_client
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetClientBySlug :one
SELECT * FROM my_client
WHERE slug = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListClients :many
SELECT * FROM my_client
WHERE deleted_at IS NULL
ORDER BY name;

-- name: CreateClient :one
INSERT INTO my_client (
    name, slug, is_project, self_capture, client_prefix,
    client_logo, address, phone_number, city, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()
) RETURNING *;

-- name: UpdateClient :one
UPDATE my_client
SET name = $2,
    slug = $3,
    is_project = $4,
    self_capture = $5,
    client_prefix = $6,
    client_logo = $7,
    address = $8,
    phone_number = $9,
    city = $10,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteClient :exec
UPDATE my_client
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;