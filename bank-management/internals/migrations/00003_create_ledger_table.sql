-- +goose Up
CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    account_id UUID NOT NULL,
    
    type VARCHAR(10) NOT NULL, 
    amount BIGINT NOT NULL CHECK (amount > 0),
    
    balance_after BIGINT NOT NULL,
    description TEXT,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),

    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

-- +goose Down
DROP TABLE ledger_entries;