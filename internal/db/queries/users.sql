-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    hashed_password
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetUser :one
SELECT user_id, uuid, name, email
FROM users
WHERE user_id = $1
LIMIT 1;

-- name: GetUserByUUID :one
SELECT *
FROM users
WHERE uuid = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: UpdateUser :one
UPDATE users
SET name = $2,
    email = $3
WHERE user_id = $1
RETURNING *;