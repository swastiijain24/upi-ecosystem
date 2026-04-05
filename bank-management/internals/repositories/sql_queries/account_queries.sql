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

-- name: GetAccountByID :one
SELECT *
FROM accounts
WHERE id = $1 AND deleted_at IS NULL;


-- name: GetBalance :one
SELECT balance
FROM accounts
WHERE id = $1 AND deleted_at IS NULL;


-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 AND deleted_at IS NULL
FOR UPDATE;


-- name: UpdateAccountBalance :exec
UPDATE accounts
SET balance = $2
WHERE id = $1;


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


-- name: CreateSettlementAccount :exec
INSERT INTO accounts (
    name,
    phone,
    is_system
)
SELECT 
    'Settlement Account',
    'SYSTEM',
    TRUE
WHERE NOT EXISTS (
    SELECT 1 
    FROM accounts 
    WHERE is_system = TRUE 
      AND name = 'Settlement Account'
);

-- name: GetSettlementAccountForUpdate :one
SELECT *
FROM accounts
WHERE name = 'Settlement Account'
  AND is_system = TRUE
  AND deleted_at IS NULL
FOR UPDATE;
