-- +goose Up
ALTER TABLE accounts
ADD CONSTRAINT accounts_balance_check
CHECK (balance >= 0 OR is_system = TRUE);

-- +goose Down
ALTER TABLE accounts
ADD CONSTRAINT accounts_balance_check
CHECK (balance >= 0);