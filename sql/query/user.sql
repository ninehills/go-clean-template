-- name: GetUser :one
SELECT * FROM user
WHERE username = ? LIMIT 1;

-- name: ListUser :many
SELECT * FROM user
ORDER BY id DESC
LIMIT ?, ?;

-- name: CreateUser :exec
INSERT INTO user (
  username, status, email, password, description, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP()
);

-- name: UpdateUser :exec
UPDATE user
SET
 status = coalesce(sqlc.narg('status'), status),
 email = coalesce(sqlc.narg('email'), email),
 password = coalesce(sqlc.narg('password'), password),
 description = coalesce(sqlc.narg('description'), description),
 updated_at = UTC_TIMESTAMP()
WHERE username = sqlc.arg('username');

-- name: DeleteUser :exec
DELETE FROM user
WHERE username = ?;