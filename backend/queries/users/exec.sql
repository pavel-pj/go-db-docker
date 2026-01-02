-- name: CreateUser :one
INSERT INTO users (email,name)
VALUES ($1, $2)
RETURNING id,name,email,created_at;

-- name: UpdateUserName :exec
UPDATE users set name = $1 where id = $2;

-- name: DeleteUser :one
DELETE from users where id = $1 
RETURNING id,email,name,created_at;