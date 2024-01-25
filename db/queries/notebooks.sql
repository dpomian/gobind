-- name: CreateNotebook :one
INSERT INTO notebooks (
  id, title, topic, content, created_at
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetNotebook :one
SELECT * FROM notebooks
WHERE id = $1 and deleted = false
LIMIT 1;


-- name: ListNotebooks :many
SELECT * FROM notebooks
WHERE deleted = false
ORDER BY title
LIMIT $1
OFFSET $2;

-- name: UpdateNotebook :one
UPDATE notebooks 
SET title = $2, content = $3, topic = $4, last_modified = $5
WHERE id=$1 and deleted = false
RETURNING *;

-- name: DeleteNotebook :exec
UPDATE notebooks 
SET deleted = true, last_modified = $2
WHERE id = $1;

-- name: SearchNotebooks :many
SELECT * from notebooks
WHERE deleted = false and (title ILIKE $1 or content ILIKE $1 or topic ILIKE $1);