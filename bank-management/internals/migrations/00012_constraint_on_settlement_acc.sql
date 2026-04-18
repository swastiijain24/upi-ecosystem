-- +goose Up
CREATE UNIQUE INDEX IF NOT EXISTS unique_system_account
ON accounts (is_system)
WHERE is_system = TRUE;

-- +goose Down
DROP INDEX IF EXISTS unique_system_account;