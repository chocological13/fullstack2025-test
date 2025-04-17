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