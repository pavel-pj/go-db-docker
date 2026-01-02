-- name: CreatePayment :one
INSERT INTO payments (invoice_id, amount_cents,status ) values($1,$2,$3)
RETURNING invoice_id,amount_cents,status,updated_at;
 
-- name: SetPaymentStatus :one 
UPDATE payments set status = $1, updated_at = CURRENT_TIMESTAMP where id = $2
RETURNING *;


