-- +goose Up
ALTER TABLE accounts
ADD COLUMN mpin_hash VARCHAR(255) ;

-- +goose Down
ALTER TABLE accounts
DROP COLUMN mpin_hash;