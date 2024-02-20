-- name: CreateFolder :exec
INSERT INTO folders (id, title, user_id) VALUES ($1, $2, $3);

-- name: GetAllByUserID :many
SELECT * FROM folders WHERE user_id = $1;

-- name: GetFolderByID :one
SELECT * FROM folders WHERE id = $1;

-- name: UpdateFolder :exec
UPDATE folders SET title = $2 WHERE id = $1;

-- name: DeleteFolder :exec
DELETE FROM folders WHERE id = $1;
