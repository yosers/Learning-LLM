-- name: InsertUserLoginOtp :exec
INSERT INTO user_login_otp (user_id, otp)
VALUES ($1, $2);

-- name: VerifyOtp :one
SELECT * FROM user_login_otp
WHERE user_id = $1 AND otp = $2 AND created_at >= NOW() - INTERVAL '5 minutes'
ORDER BY created_at DESC
LIMIT 1;