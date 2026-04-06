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

-- name: BalanceFromEntries :one 
SELECT CAST(
    COALESCE(SUM(credit), 0) - COALESCE(SUM(debit), 0)
    AS BIGINT
) AS calculated_balance
FROM ledger_entries
WHERE account_id = $1;