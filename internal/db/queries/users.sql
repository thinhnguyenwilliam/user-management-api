-- name: CreateUser :one
INSERT INTO users (
    user_fullname,
    user_email,
    user_password,
    user_age,
    user_status,
    user_level
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;


-- name: GetUserByUUID :one
SELECT *
FROM users
WHERE user_uuid = $1
AND user_deleted_at IS NULL;


-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE user_email = $1
AND user_deleted_at IS NULL;


-- name: ListUsers :many
SELECT *
FROM users
WHERE user_deleted_at IS NULL
ORDER BY user_created_at DESC
LIMIT $1
OFFSET $2;


-- name: UpdateUser :one
UPDATE users
SET
    user_fullname = $2,
    user_age = $3,
    user_status = $4,
    user_level = $5
WHERE user_uuid = $1
AND user_deleted_at IS NULL
RETURNING *;


-- name: UpdateUserPassword :exec
UPDATE users
SET user_password = $2
WHERE user_uuid = $1;


-- name: DeleteUserSoft :exec
UPDATE users
SET user_deleted_at = NOW()
WHERE user_uuid = $1;


-- name: DeleteUserHard :exec
DELETE FROM users
WHERE user_uuid = $1;