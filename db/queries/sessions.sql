-- name: CreateSession :one
INSERT INTO sessions (
    user_id,
    channel_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetCurrentSessions :one
SELECT DISTINCT s.*
FROM sessions s
JOIN conversations c ON c.session_id = s.id
WHERE s.channel_id = $1
  AND c.created_at >= NOW() - INTERVAL '1 day' AND s.user_id = $2
LIMIT 1;
