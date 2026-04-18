-- +goose Up
ALTER TABLE transactions 
ADD CONSTRAINT unique_external_id_type UNIQUE (external_id, type);

-- +goose Down
ALTER TABLE transactions 
DROP CONSTRAINT IF EXISTS unique_external_id_type;
