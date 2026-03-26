-- +goose Up
ALTER TABLE accounts
ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE accounts
DROP COLUMN is_system;