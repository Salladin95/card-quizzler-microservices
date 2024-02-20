-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS folders (
    id VARCHAR(56) PRIMARY KEY,
    title VARCHAR(56) NOT NULL,
    user_id VARCHAR(56) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS modules (
    id VARCHAR(56) PRIMARY KEY,
    title VARCHAR(56) NOT NULL,
    user_id VARCHAR(56) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS module_folder (
    module_id VARCHAR(56),
    folder_id VARCHAR(56),
    PRIMARY KEY (module_id, folder_id),
    FOREIGN KEY (module_id) REFERENCES modules(id),
    FOREIGN KEY (folder_id) REFERENCES folders(id)
);

CREATE TABLE IF NOT EXISTS terms (
    id VARCHAR(56) PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    module_id VARCHAR(56) NOT NULL REFERENCES modules(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS terms;
DROP TABLE IF EXISTS module_folder;
DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS folders;
-- +goose StatementEnd
