-- name: CreateModule :exec
INSERT INTO modules (id, title, user_id) VALUES ($1, $2, $3);

-- name: GetAllModulesByUserID :many
SELECT
       m.id AS module_id,
       m.title AS module_title
FROM modules m
WHERE m.user_id = $1;

-- name: GetModuleByID :one
SELECT
    m.id AS module_id,
    m.title AS module_title,
    m.user_id AS user_id,
    json_agg(json_build_object(
            'id', t.id,
            'title', t.title,
            'description', t.description
             )) AS terms
FROM
    modules m
        JOIN
    terms t ON m.id = t.module_id
WHERE
    m.id = $1
GROUP BY
    m.id, m.title;

-- name: UpdateModule :exec
UPDATE modules SET title = $2 WHERE id = $1;

-- name: DeleteModule :exec
DELETE FROM modules WHERE id = $1;
