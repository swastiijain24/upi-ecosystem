-- +goose Up

-- drop FK first
ALTER TABLE transactions
DROP CONSTRAINT IF EXISTS transactions_from_account_id_fkey;

-- change type (safe since no data)
ALTER TABLE transactions
ALTER COLUMN from_account_id TYPE TEXT;

-- +goose Down

ALTER TABLE transactions
ALTER COLUMN from_account_id TYPE UUID;
