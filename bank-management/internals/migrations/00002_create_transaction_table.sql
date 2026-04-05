-- +goose Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    from_account_id TEXT NOT NULL,  -- changed directly to TEXT
    to_account_identifier TEXT NOT NULL,

    amount BIGINT NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS transactions;