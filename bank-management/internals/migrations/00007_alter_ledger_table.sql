-- +goose Up
ALTER TABLE ledger_entries
ADD COLUMN transaction_id UUID NOT NULL;

ALTER TABLE ledger_entries
ADD CONSTRAINT fk_ledger_transaction
FOREIGN KEY (transaction_id)
REFERENCES transactions(id);

-- +goose Down
ALTER TABLE ledger_entries
DROP CONSTRAINT fk_ledger_transaction;

ALTER TABLE ledger_entries
DROP COLUMN transaction_id;