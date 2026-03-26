-- name: GetAccountByID :one
SELECT *
FROM accounts
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateAccount :one
INSERT INTO accounts (
    name,
    phone
)
VALUES (
    $1,
    $2
)
RETURNING id, name, phone, balance, created_at;

-- name: GetBalance :one
SELECT balance
FROM accounts
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 AND deleted_at IS NULL
FOR UPDATE;

-- name: CheckIdempotencyKey :one 
SELECT * FROM idempotency_keys WHERE key = $1;

-- name: CreateLedgerEntry :exec
INSERT INTO ledger_entries (
    account_id,
    type,
    amount,
    balance_after,
    description
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: CreateTransaction :one
INSERT INTO transactions (
    from_account_id,
    to_account_identifier,
    amount,
    status
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, from_account_id, to_account_identifier, amount, status, created_at, updated_at;


-- name: UpdateAccountBalance :exec
UPDATE accounts
SET balance = $2
WHERE id = $1;


-- name: GetTransactions :many
SELECT *
FROM transactions
WHERE from_account_id = $1 or to_account_identifier = $1
ORDER BY created_at DESC;

-- name: DeleteAccount :exec
UPDATE accounts
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UpdatePaymentStatus :exec
UPDATE transactions
SET status = $2,
    updated_at = NOW()
WHERE id = $1
  AND status != $2;