-- name: CreateProduct :one
INSERT INTO products (name, price,status)
VALUES ($1, $2,$3)
RETURNING id, name, price,status, created_at, updated_at;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: GetProducts :many
SELECT * FROM products;