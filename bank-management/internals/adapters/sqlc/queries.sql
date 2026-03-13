-- name: GetAccountByID :one
SELECT *
FROM accounts
WHERE id = $1;

-- name: CreateAccount :one
INSERT INTO accounts (
    id,
    name,
    phone,
    balance
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, name, phone, balance, created_at;

-- name: UpdateAccountBalance :exec
UPDATE accounts
SET balance = $2
WHERE id = $1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1
FOR UPDATE;

-- name: CreateTransaction :one
INSERT INTO transactions (
    id,
    account_id,
    type,
    amount,
    status,
    created_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    NOW()
)
RETURNING *;

-- name: GetTransactions :many
SELECT 
    id,
    account_id,
    type,
    amount,
    status,
    created_at
FROM transactions
WHERE account_id = $1
ORDER BY created_at DESC;

-- name: CheckBalance :one
SELECT balance
FROM accounts
WHERE id = $1;