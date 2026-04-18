-- name: CreateAccount :one
INSERT INTO accounts (
    name,
    phone,
    mpin_hash
)
VALUES (
    $1,
    $2,
    $3
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


-- name: UpdateUserBalanceDebit :one
UPDATE accounts 
SET balance = balance - $1, updated_at = NOW() 
WHERE id = $2 
  AND balance >= $1  
  AND deleted_at IS NULL
RETURNING balance;

-- name: UpdateAccountBalanceCredit :one
UPDATE accounts 
SET balance = balance + $1, updated_at = NOW() 
WHERE id = $2 
  AND deleted_at IS NULL
RETURNING balance;

-- name: UpdateSettlementBalanceAtomic :one
UPDATE accounts 
SET balance = balance + $1, updated_at = NOW() 
WHERE id = $2 AND deleted_at IS NULL
RETURNING balance;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 AND deleted_at IS NULL
FOR UPDATE;

-- name: DeleteAccount :exec
UPDATE accounts
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: CreateSettlementAccount :one
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
)
RETURNING id;

-- name: GetSettlementAccount :one
SELECT *
FROM accounts
WHERE name = 'Settlement Account'
  AND is_system = TRUE
  AND deleted_at IS NULL;


-- name: GetMpinHash :one
SELECT mpin_hash 
FROM accounts 
WHERE id = $1 
LIMIT 1;