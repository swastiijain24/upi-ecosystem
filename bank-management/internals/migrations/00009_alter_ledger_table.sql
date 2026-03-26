-- +goose Up
ALTER TABLE ledger_entries
DROP COLUMN amount;

ALTER TABLE ledger_entries
ADD COLUMN debit BIGINT NOT NULL DEFAULT 0 CHECK (debit >= 0),
ADD COLUMN credit BIGINT NOT NULL DEFAULT 0 CHECK (credit >= 0);


-- +goose Down
ALTER TABLE ledger_entries
DROP COLUMN debit,
DROP COLUMN credit;

ALTER TABLE ledger_entries
ADD COLUMN amount BIGINT NOT NULL DEFAULT 0;