-- name: CreateTerm :exec
INSERT INTO terms (id, title, description, module_id) VALUES ($1, $2, $3, $4);

-- name: GetTermsByModuleID :many
SELECT * FROM terms WHERE module_id = $1;

-- name: GetTermByID :one
SELECT * FROM terms WHERE id = $1;

-- name: UpdateTerm :exec
UPDATE terms SET title = $2, description = $3 WHERE id = $1;

-- name: DeleteTerm :exec
DELETE FROM terms WHERE id = $1;
