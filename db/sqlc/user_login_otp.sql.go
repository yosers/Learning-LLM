// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user_login_otp.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countValidOtps = `-- name: CountValidOtps :one
SELECT COUNT(*) FROM user_login_otp
WHERE otp = $1
  AND expires_at >= (NOW() AT TIME ZONE 'UTC')
`

func (q *Queries) CountValidOtps(ctx context.Context, otp string) (int64, error) {
	row := q.db.QueryRow(ctx, countValidOtps, otp)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const findUserLoginOtpByPhone = `-- name: FindUserLoginOtpByPhone :one
SELECT ulo.id, ulo.user_id, ulo.otp, ulo.is_used, ulo.created_at, ulo.expires_at FROM user_login_otp ulo join users u
ON ulo.user_id = u.id
WHERE u.phone = $1
ORDER BY ulo. created_at DESC
LIMIT 1
`

func (q *Queries) FindUserLoginOtpByPhone(ctx context.Context, phone pgtype.Text) (UserLoginOtp, error) {
	row := q.db.QueryRow(ctx, findUserLoginOtpByPhone, phone)
	var i UserLoginOtp
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Otp,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const findUserLoginOtpNotActive = `-- name: FindUserLoginOtpNotActive :one
SELECT id, user_id, otp, is_used, created_at, expires_at FROM user_login_otp ulo
WHERE ulo.user_id = $1 AND is_used = FALSE
`

func (q *Queries) FindUserLoginOtpNotActive(ctx context.Context, userID int32) (UserLoginOtp, error) {
	row := q.db.QueryRow(ctx, findUserLoginOtpNotActive, userID)
	var i UserLoginOtp
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Otp,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const insertUserLoginOtp = `-- name: InsertUserLoginOtp :exec
INSERT INTO user_login_otp (user_id, otp)
VALUES ($1, $2)
`

type InsertUserLoginOtpParams struct {
	UserID int32
	Otp    string
}

func (q *Queries) InsertUserLoginOtp(ctx context.Context, arg InsertUserLoginOtpParams) error {
	_, err := q.db.Exec(ctx, insertUserLoginOtp, arg.UserID, arg.Otp)
	return err
}

const updateIsUsed = `-- name: UpdateIsUsed :exec
UPDATE user_login_otp
SET is_used = TRUE 
WHERE user_id = $1 AND otp = $2
`

type UpdateIsUsedParams struct {
	UserID int32
	Otp    string
}

func (q *Queries) UpdateIsUsed(ctx context.Context, arg UpdateIsUsedParams) error {
	_, err := q.db.Exec(ctx, updateIsUsed, arg.UserID, arg.Otp)
	return err
}

const updateIsUsedFalse = `-- name: UpdateIsUsedFalse :exec
UPDATE user_login_otp
SET is_used = FALSE
WHERE user_id = $1
`

func (q *Queries) UpdateIsUsedFalse(ctx context.Context, userID int32) error {
	_, err := q.db.Exec(ctx, updateIsUsedFalse, userID)
	return err
}

const updateOTPByUserId = `-- name: UpdateOTPByUserId :exec
UPDATE user_login_otp
SET is_used = FALSE, otp = $1, created_at = NOW(), 
    expires_at = NOW() + INTERVAL '5 minutes'
WHERE user_id = $2
`

type UpdateOTPByUserIdParams struct {
	Otp    string
	UserID int32
}

func (q *Queries) UpdateOTPByUserId(ctx context.Context, arg UpdateOTPByUserIdParams) error {
	_, err := q.db.Exec(ctx, updateOTPByUserId, arg.Otp, arg.UserID)
	return err
}

const verifyOtp = `-- name: VerifyOtp :one
SELECT CASE
       WHEN expires_at < (NOW() AT TIME ZONE 'UTC') THEN 'EXPIRED'
       WHEN is_used = TRUE THEN 'USED'
       ELSE 'VALID'
     END as status, user_id
FROM user_login_otp
WHERE otp = $1
ORDER BY created_at DESC
LIMIT 1
`

type VerifyOtpRow struct {
	Status string
	UserID int32
}

func (q *Queries) VerifyOtp(ctx context.Context, otp string) (VerifyOtpRow, error) {
	row := q.db.QueryRow(ctx, verifyOtp, otp)
	var i VerifyOtpRow
	err := row.Scan(&i.Status, &i.UserID)
	return i, err
}
