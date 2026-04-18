-- +goose Up
ALTER TABLE transactions 
ADD COLUMN type VARCHAR(20) NOT NULL;

-- +goose Down
ALTER TABLE transactions 
DROP COLUMN type;
