-- +goose Up
CREATE TYPE transaction_state AS ENUM (
    'INITIATED',
    'DEBIT_PENDING',
    'DEBIT_SUCCESS',
    'DEBIT_FAILED',
    'CREDIT_PENDING',
    'CREDIT_SUCCESS',
    'CREDIT_FAILED',
    'SUCCESS'
    'FAILED',
    'REVERSED'
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,

    payer_vpa TEXT NOT NULL,
    payer_bank TEXT NOT NULL,

    payee_vpa TEXT NOT NULL,
    payee_bank TEXT NOT NULL,

    amount BIGINT NOT NULL,

    state transaction_state NOT NULL,

    reference_id TEXT UNIQUE,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE transactions;
DROP TYPE transaction_state;