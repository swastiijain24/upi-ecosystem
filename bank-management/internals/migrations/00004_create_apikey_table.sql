-- +goose Up
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    key_id TEXT NOT NULL UNIQUE,
    key_hash TEXT NOT NULL,

    client TEXT NOT NULL,                

    allowed_ips TEXT[] DEFAULT '{}',

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,

    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_api_keys_key_id ON api_keys(key_id);


-- +goose Down
DROP TABLE IF EXISTS api_keys;

