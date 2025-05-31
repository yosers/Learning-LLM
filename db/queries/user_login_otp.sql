-- name: InsertUserLoginOtp :exec
INSERT INTO user_login_otp (user_id, otp)
VALUES ($1, $2);

-- name: VerifyOtp :one
SELECT * FROM user_login_otp
WHERE user_id = $1 AND otp = $2 AND created_at >= NOW() - INTERVAL '5 minutes' AND is_used = FALSE
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateIsUsed :exec
UPDATE user_login_otp
SET is_used = TRUE 
WHERE user_id = $1 AND otp = $2;

-- name: FindUserLoginOtpByPhone :one
SELECT * FROM user_login_otp ulo join users u
ON ulo.user_id = u.id
WHERE u.phone = $1 AND is_used = FALSE
ORDER BY ulo. created_at DESC
LIMIT 1;

-- name: FindUserLoginOtpNotActive :one
SELECT * FROM user_login_otp ulo
WHERE ulo.user_id = $1 AND is_used = FALSE;