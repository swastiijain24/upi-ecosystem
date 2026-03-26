-- +goose Up
CREATE TABLE idempotency_keys (
    key VARCHAR(255) PRIMARY KEY,
    response_code INT,
    response_body JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE idempotency_keys;