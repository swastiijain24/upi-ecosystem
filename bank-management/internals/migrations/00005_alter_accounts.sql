-- +goose Up
ALTER TABLE accounts
ADD COLUMN phone VARCHAR(20) NOT NULL;

-- +goose Down
ALTER TABLE accounts
DROP COLUMN phone;

