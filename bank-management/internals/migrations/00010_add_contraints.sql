-- +goose Up

ALTER TABLE ledger_entries
ADD CONSTRAINT check_single_side
CHECK (
    (debit > 0 AND credit = 0) OR 
    (debit = 0 AND credit > 0)
);


-- +goose Down

ALTER TABLE ledger_entries
DROP CONSTRAINT check_single_side;