
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


-- name: GetTransactions :many
SELECT *
FROM transactions
WHERE from_account_id = $1 or to_account_identifier = $1
ORDER BY created_at DESC;

-- name: GetTransactionById :one 
SELECT * 
FROM transactions 
WHERE ID = $1 ;

-- name: UpdatePaymentStatus :exec
UPDATE transactions
SET status = $2,
    updated_at = NOW()
WHERE id = $1
  AND status != $2;
