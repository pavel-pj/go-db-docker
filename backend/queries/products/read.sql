-- name: GetProductByID :one
SELECT id,slug,title,description,price_cents,created_at 
from products where id = $1;

-- name: ListProducts :many 
SELECT id,slug,title,description,price_cents,created_at FROM products
LIMIT $1 OFFSET $2;