-- +goose Up
ALTER TABLE accounts
DROP CONSTRAINT IF EXISTS accounts_balance_check;

-- +goose Down
ALTER TABLE accounts
ADD CONSTRAINT accounts_balance_check
CHECK (balance >= 0);