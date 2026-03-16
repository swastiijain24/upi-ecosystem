-- name: CreateTransaction :one
INSERT INTO transactions (
    id,
    payer_vpa,
    payer_bank,
    payee_vpa,
    payee_bank,
    amount,
    state,
    reference_id
) VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8
)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1;

-- name: UpdateTransactionState :one
UPDATE transactions
SET state = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetTransactionByReference :one
SELECT * FROM transactions
WHERE reference_id = $1;