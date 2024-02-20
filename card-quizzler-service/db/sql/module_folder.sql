-- name: GetFolderWithModulesForUser :many
SELECT
    f.id AS folder_id,
    f.title AS folder_title,
    sqlc.embed(m)
FROM folders f
LEFT JOIN module_folder mf ON f.id = mf.folder_id
LEFT JOIN modules m ON mf.module_id = m.id
WHERE f.user_id = $1;


-- name: AddModuleToFolder :exec
INSERT INTO module_folder (module_id, folder_id) VALUES ($1, $2);
