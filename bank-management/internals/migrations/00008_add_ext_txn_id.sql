-- +goose Up
ALTER TABLE transactions
ADD COLUMN external_id TEXT NOT NULL;

-- +goose Down
ALTER TABLE transactions
DROP COLUMN external_id;