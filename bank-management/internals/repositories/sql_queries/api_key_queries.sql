-- name: CreateAPIKey :one
INSERT INTO api_keys (
    key_id,
    key_hash,
    client,
    allowed_ips,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING 
    id,
    key_id,
    key_hash,
    client,
    allowed_ips,
    created_at,
    expires_at,
    last_used_at,
    is_active;


-- name: GetAPIKeyByKeyID :one
SELECT 
    id,
    key_id,
    key_hash,
    client,
    allowed_ips,
    created_at,
    expires_at,
    last_used_at,
    is_active
FROM api_keys
WHERE key_id = $1
LIMIT 1;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE api_keys
SET last_used_at = NOW()
WHERE key_id = $1;

-- name: DeactivateAPIKey :exec
UPDATE api_keys
SET is_active = FALSE
WHERE key_id = $1;


-- name: IsValid :one
SELECT EXISTS (
    SELECT 1 FROM api_keys
    WHERE key_id = $1
      AND is_active = TRUE
      AND (expires_at IS NULL OR expires_at > NOW())
);