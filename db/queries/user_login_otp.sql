-- name: InsertUserLoginOtp :exec
INSERT INTO user_login_otp (user_id, otp)
VALUES ($1, $2);

-- name: VerifyOtp :one
SELECT * FROM user_login_otp
WHERE user_id = $1 
    AND otp = $2
    AND expires_at >= (NOW() AT TIME ZONE 'UTC')
    AND is_used = FALSE
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateIsUsed :exec
UPDATE user_login_otp
SET is_used = TRUE 
WHERE user_id = $1 AND otp = $2;

-- name: FindUserLoginOtpByPhone :one
SELECT ulo.* FROM user_login_otp ulo join users u
ON ulo.user_id = u.id
WHERE u.phone = $1
ORDER BY ulo. created_at DESC
LIMIT 1;

-- name: FindUserLoginOtpNotActive :one
SELECT * FROM user_login_otp ulo
WHERE ulo.user_id = $1 AND is_used = FALSE;

-- name: UpdateIsUsedFalse :exec
UPDATE user_login_otp
SET is_used = FALSE
WHERE user_id = $1;

-- name: UpdateOTPByUserId :exec
UPDATE user_login_otp
SET is_used = FALSE, otp = $1, created_at = NOW(), 
    expires_at = NOW() + INTERVAL '5 minutes'
WHERE user_id = $2;