-- name: CreateProduct :one
INSERT INTO products (name, status, price)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteProduct :execrows
DELETE from products where name = $1;
 