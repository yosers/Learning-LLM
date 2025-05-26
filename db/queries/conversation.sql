-- name: CreateConversation :one
INSERT INTO conversations (
    session_id,
    message,
    role
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetConversationsBySessionID :many
SELECT id, session_id, message, role, created_at, updated_at
FROM conversations
WHERE session_id = $1
ORDER BY created_at ASC; 