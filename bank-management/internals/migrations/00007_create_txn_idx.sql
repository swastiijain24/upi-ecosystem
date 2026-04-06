-- +goose Up
CREATE INDEX idx_transactions_from_account 
ON transactions(from_account_id);

CREATE INDEX idx_transactions_to_account 
ON transactions(to_account_identifier);

-- +goose Down
DROP INDEX IF EXISTS idx_transactions_from_account;
DROP INDEX IF EXISTS idx_transactions_to_account;