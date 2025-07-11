// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: shops.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countShops = `-- name: CountShops :one
SELECT COUNT(*) FROM shops
WHERE is_active = true
`

func (q *Queries) CountShops(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countShops)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createShops = `-- name: CreateShops :one
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
RETURNING id, name, description, logo_url, website_url, email, whatsapp_phone, address, city, state, zip_code, country, latitude, longitude, is_active, slug, created_at, updated_at
`

type CreateShopsParams struct {
	Name          string
	Description   string
	LogoUrl       pgtype.Text
	WebsiteUrl    pgtype.Text
	Email         pgtype.Text
	WhatsappPhone pgtype.Text
	Address       string
	City          string
	State         string
	ZipCode       string
	Country       string
	Latitude      float64
	Longitude     float64
	IsActive      bool
}

func (q *Queries) CreateShops(ctx context.Context, arg CreateShopsParams) (Shop, error) {
	row := q.db.QueryRow(ctx, createShops,
		arg.Name,
		arg.Description,
		arg.LogoUrl,
		arg.WebsiteUrl,
		arg.Email,
		arg.WhatsappPhone,
		arg.Address,
		arg.City,
		arg.State,
		arg.ZipCode,
		arg.Country,
		arg.Latitude,
		arg.Longitude,
		arg.IsActive,
	)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.LogoUrl,
		&i.WebsiteUrl,
		&i.Email,
		&i.WhatsappPhone,
		&i.Address,
		&i.City,
		&i.State,
		&i.ZipCode,
		&i.Country,
		&i.Latitude,
		&i.Longitude,
		&i.IsActive,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteShopsById = `-- name: DeleteShopsById :exec
UPDATE shops    
SET is_active = false, updated_at = now()
WHERE id = $1
`

func (q *Queries) DeleteShopsById(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteShopsById, id)
	return err
}

const getAllShops = `-- name: GetAllShops :many
SELECT s.id, s.name, s.description, s.logo_url, s.website_url, s.email, s.whatsapp_phone, s.address, s.city, s.state, s.zip_code, s.country, s.latitude, s.longitude, s.is_active, s.slug, s.created_at, s.updated_at FROM shops s 
WHERE  s.is_active = true
ORDER BY s.created_at DESC
`

func (q *Queries) GetAllShops(ctx context.Context) ([]Shop, error) {
	rows, err := q.db.Query(ctx, getAllShops)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Shop
	for rows.Next() {
		var i Shop
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.LogoUrl,
			&i.WebsiteUrl,
			&i.Email,
			&i.WhatsappPhone,
			&i.Address,
			&i.City,
			&i.State,
			&i.ZipCode,
			&i.Country,
			&i.Latitude,
			&i.Longitude,
			&i.IsActive,
			&i.Slug,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getShopsById = `-- name: GetShopsById :one
SELECT s.id, s.name, s.description, s.logo_url, s.website_url, s.email, s.whatsapp_phone, s.address, s.city, s.state, s.zip_code, s.country, s.latitude, s.longitude, s.is_active, s.slug, s.created_at, s.updated_at FROM shops s 
WHERE s.id = $1 and s.is_active = true LIMIT 1
`

func (q *Queries) GetShopsById(ctx context.Context, id int32) (Shop, error) {
	row := q.db.QueryRow(ctx, getShopsById, id)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.LogoUrl,
		&i.WebsiteUrl,
		&i.Email,
		&i.WhatsappPhone,
		&i.Address,
		&i.City,
		&i.State,
		&i.ZipCode,
		&i.Country,
		&i.Latitude,
		&i.Longitude,
		&i.IsActive,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getShopsByNameOrWhatshapp = `-- name: GetShopsByNameOrWhatshapp :one
SELECT s.id, s.name, s.description, s.logo_url, s.website_url, s.email, s.whatsapp_phone, s.address, s.city, s.state, s.zip_code, s.country, s.latitude, s.longitude, s.is_active, s.slug, s.created_at, s.updated_at FROM shops s 
WHERE s.is_active = true 
  AND (
    s.name = $1 
    OR s.whatsapp_phone = $2
) LIMIT 1
`

type GetShopsByNameOrWhatshappParams struct {
	Name          string
	WhatsappPhone pgtype.Text
}

func (q *Queries) GetShopsByNameOrWhatshapp(ctx context.Context, arg GetShopsByNameOrWhatshappParams) (Shop, error) {
	row := q.db.QueryRow(ctx, getShopsByNameOrWhatshapp, arg.Name, arg.WhatsappPhone)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.LogoUrl,
		&i.WebsiteUrl,
		&i.Email,
		&i.WhatsappPhone,
		&i.Address,
		&i.City,
		&i.State,
		&i.ZipCode,
		&i.Country,
		&i.Latitude,
		&i.Longitude,
		&i.IsActive,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listShops = `-- name: ListShops :many
SELECT s.id, s.name, s.description, s.logo_url, s.website_url, s.email, s.whatsapp_phone, s.address, s.city, s.state, s.zip_code, s.country, s.latitude, s.longitude, s.is_active, s.slug, s.created_at, s.updated_at FROM shops s 
WHERE  s.is_active = true
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2
`

type ListShopsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListShops(ctx context.Context, arg ListShopsParams) ([]Shop, error) {
	rows, err := q.db.Query(ctx, listShops, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Shop
	for rows.Next() {
		var i Shop
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.LogoUrl,
			&i.WebsiteUrl,
			&i.Email,
			&i.WhatsappPhone,
			&i.Address,
			&i.City,
			&i.State,
			&i.ZipCode,
			&i.Country,
			&i.Latitude,
			&i.Longitude,
			&i.IsActive,
			&i.Slug,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateShops = `-- name: UpdateShops :one
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
RETURNING id, name, description, logo_url, website_url, email, whatsapp_phone, address, city, state, zip_code, country, latitude, longitude, is_active, slug, created_at, updated_at
`

type UpdateShopsParams struct {
	ID            int32
	Name          string
	Description   string
	LogoUrl       pgtype.Text
	WebsiteUrl    pgtype.Text
	Email         pgtype.Text
	WhatsappPhone pgtype.Text
	Address       string
	City          string
	State         string
	ZipCode       string
	Country       string
	Latitude      float64
	Longitude     float64
	IsActive      bool
}

func (q *Queries) UpdateShops(ctx context.Context, arg UpdateShopsParams) (Shop, error) {
	row := q.db.QueryRow(ctx, updateShops,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.LogoUrl,
		arg.WebsiteUrl,
		arg.Email,
		arg.WhatsappPhone,
		arg.Address,
		arg.City,
		arg.State,
		arg.ZipCode,
		arg.Country,
		arg.Latitude,
		arg.Longitude,
		arg.IsActive,
	)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.LogoUrl,
		&i.WebsiteUrl,
		&i.Email,
		&i.WhatsappPhone,
		&i.Address,
		&i.City,
		&i.State,
		&i.ZipCode,
		&i.Country,
		&i.Latitude,
		&i.Longitude,
		&i.IsActive,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
