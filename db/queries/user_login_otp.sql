-- name: InsertUserLoginOtp :exec
INSERT INTO user_login_otp (user_id, otp)
VALUES ($1, $2);