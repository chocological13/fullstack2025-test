// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"
)

const createClient = `-- name: CreateClient :one
INSERT INTO my_client (
    name, slug, is_project, self_capture, client_prefix,
    client_logo, address, phone_number, city, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()
) RETURNING id, name, slug, is_project, self_capture, client_prefix, client_logo, address, phone_number, city, created_at, updated_at, deleted_at
`

type CreateClientParams struct {
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	IsProject    string         `json:"is_project"`
	SelfCapture  string         `json:"self_capture"`
	ClientPrefix string         `json:"client_prefix"`
	ClientLogo   string         `json:"client_logo"`
	Address      sql.NullString `json:"address"`
	PhoneNumber  sql.NullString `json:"phone_number"`
	City         sql.NullString `json:"city"`
}

func (q *Queries) CreateClient(ctx context.Context, arg CreateClientParams) (MyClient, error) {
	row := q.queryRow(ctx, q.createClientStmt, createClient,
		arg.Name,
		arg.Slug,
		arg.IsProject,
		arg.SelfCapture,
		arg.ClientPrefix,
		arg.ClientLogo,
		arg.Address,
		arg.PhoneNumber,
		arg.City,
	)
	var i MyClient
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.IsProject,
		&i.SelfCapture,
		&i.ClientPrefix,
		&i.ClientLogo,
		&i.Address,
		&i.PhoneNumber,
		&i.City,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteClient = `-- name: DeleteClient :exec
UPDATE my_client
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) DeleteClient(ctx context.Context, id int32) error {
	_, err := q.exec(ctx, q.deleteClientStmt, deleteClient, id)
	return err
}

const getClient = `-- name: GetClient :one
SELECT id, name, slug, is_project, self_capture, client_prefix, client_logo, address, phone_number, city, created_at, updated_at, deleted_at FROM my_client
WHERE id = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetClient(ctx context.Context, id int32) (MyClient, error) {
	row := q.queryRow(ctx, q.getClientStmt, getClient, id)
	var i MyClient
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.IsProject,
		&i.SelfCapture,
		&i.ClientPrefix,
		&i.ClientLogo,
		&i.Address,
		&i.PhoneNumber,
		&i.City,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getClientBySlug = `-- name: GetClientBySlug :one
SELECT id, name, slug, is_project, self_capture, client_prefix, client_logo, address, phone_number, city, created_at, updated_at, deleted_at FROM my_client
WHERE slug = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetClientBySlug(ctx context.Context, slug string) (MyClient, error) {
	row := q.queryRow(ctx, q.getClientBySlugStmt, getClientBySlug, slug)
	var i MyClient
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.IsProject,
		&i.SelfCapture,
		&i.ClientPrefix,
		&i.ClientLogo,
		&i.Address,
		&i.PhoneNumber,
		&i.City,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listClients = `-- name: ListClients :many
SELECT id, name, slug, is_project, self_capture, client_prefix, client_logo, address, phone_number, city, created_at, updated_at, deleted_at FROM my_client
WHERE deleted_at IS NULL
ORDER BY name
`

func (q *Queries) ListClients(ctx context.Context) ([]MyClient, error) {
	rows, err := q.query(ctx, q.listClientsStmt, listClients)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MyClient
	for rows.Next() {
		var i MyClient
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.IsProject,
			&i.SelfCapture,
			&i.ClientPrefix,
			&i.ClientLogo,
			&i.Address,
			&i.PhoneNumber,
			&i.City,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateClient = `-- name: UpdateClient :one
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
RETURNING id, name, slug, is_project, self_capture, client_prefix, client_logo, address, phone_number, city, created_at, updated_at, deleted_at
`

type UpdateClientParams struct {
	ID           int32          `json:"id"`
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	IsProject    string         `json:"is_project"`
	SelfCapture  string         `json:"self_capture"`
	ClientPrefix string         `json:"client_prefix"`
	ClientLogo   string         `json:"client_logo"`
	Address      sql.NullString `json:"address"`
	PhoneNumber  sql.NullString `json:"phone_number"`
	City         sql.NullString `json:"city"`
}

func (q *Queries) UpdateClient(ctx context.Context, arg UpdateClientParams) (MyClient, error) {
	row := q.queryRow(ctx, q.updateClientStmt, updateClient,
		arg.ID,
		arg.Name,
		arg.Slug,
		arg.IsProject,
		arg.SelfCapture,
		arg.ClientPrefix,
		arg.ClientLogo,
		arg.Address,
		arg.PhoneNumber,
		arg.City,
	)
	var i MyClient
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.IsProject,
		&i.SelfCapture,
		&i.ClientPrefix,
		&i.ClientLogo,
		&i.Address,
		&i.PhoneNumber,
		&i.City,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
