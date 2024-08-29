-- name: RegisterUser :one
INSERT INTO
    users (
    email,
    hashed_password,
    first_name,
    last_name,
    user_role
  )
VALUES (
  $1, $2, $3, $4, $5
)
RETURNING
  *;

-- name: GetUser :one
SELECT
  *
FROM
  users
WHERE
  email = $1
LIMIT
1;