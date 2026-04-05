-- name: CreateLedgerEntry :exec
INSERT INTO ledger_entries (
    transaction_id,
    account_id,
    type,
    debit,
    credit,
    balance_after,
    description
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);