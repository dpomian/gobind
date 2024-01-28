-- name: CreateNotebook :one
INSERT INTO notebooks (
  id, title, topic, content, user_id, created_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetNotebook :one
SELECT * FROM notebooks
WHERE id = $1 and user_id = $2 and deleted = false
LIMIT 1;


-- name: ListNotebooks :many
SELECT * FROM notebooks
WHERE user_id = $1 and deleted = false
ORDER BY title
LIMIT $2
OFFSET $3;

-- name: UpdateNotebook :one
UPDATE notebooks 
SET title = $3, content = $4, topic = $5, last_modified = $6
WHERE id=$1 and user_id = $2 and deleted = false 
RETURNING *;

-- name: DeleteNotebook :one
UPDATE notebooks 
SET deleted = true, last_modified = $3
WHERE id = $1 and user_id = $2
RETURNING *;

-- name: SearchNotebooks :many
SELECT * from notebooks
WHERE user_id = $1 and deleted = false and (title ILIKE $2 or content ILIKE $2 or topic ILIKE $2);