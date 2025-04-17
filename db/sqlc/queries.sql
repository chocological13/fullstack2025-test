-- name: GetClient :one
SELECT * FROM my_client
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;