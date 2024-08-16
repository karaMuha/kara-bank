-- name: RegisterUser :one
INSERT INTO
    users (
    email,
    hashed_password,
    first_name,
    last_name
  )
VALUES (
  $1, $2, $3, $4
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