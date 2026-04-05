-- +goose Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_id UUID NOT NULL,
    transaction_id UUID NOT NULL,

    type VARCHAR(10) NOT NULL,

    debit BIGINT NOT NULL DEFAULT 0 CHECK (debit >= 0),
    credit BIGINT NOT NULL DEFAULT 0 CHECK (credit >= 0),

    balance_after BIGINT NOT NULL,
    description TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT fk_ledger_account
        FOREIGN KEY (account_id) REFERENCES accounts(id),

    CONSTRAINT fk_ledger_transaction
        FOREIGN KEY (transaction_id) REFERENCES transactions(id),

    CONSTRAINT check_single_side
        CHECK (
            (debit > 0 AND credit = 0) OR 
            (debit = 0 AND credit > 0)
        )
);

CREATE INDEX idx_ledger_account_id ON ledger_entries(account_id);
CREATE INDEX idx_ledger_transaction_id ON ledger_entries(transaction_id);

-- +goose Down
DROP TABLE IF EXISTS ledger_entries;