-- +goose Up
-- +goose StatementBegin

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create the folders table
CREATE TABLE IF NOT EXISTS folders (
    id UUID PRIMARY KEY,
    title VARCHAR(255) UNIQUE NOT NULL,
    user_id VARCHAR(255) REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_open BOOLEAN DEFAULT FALSE,
    copies_coint INT DEFAULT 0
);

-- Create the modules table
CREATE TABLE IF NOT EXISTS modules (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_open BOOLEAN DEFAULT FALSE,
    copies_coint INT DEFAULT 0
);

-- Create the module_folders junction table
CREATE TABLE IF NOT EXISTS module_folders (
    module_id UUID REFERENCES modules(id),
    folder_id UUID REFERENCES folders(id),
    PRIMARY KEY (module_id, folder_id)
);

-- Create the terms table
CREATE TABLE IF NOT EXISTS terms (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    negative_answer_streak INT NOT NULL DEFAULT 0,
    positive_answer_streak INT NOT NULL DEFAULT 0,
    is_difficult BOOLEAN NOT NULL DEFAULT FALSE,
    module_id UUID REFERENCES modules(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS module_terms;
DROP TABLE IF EXISTS module_folders;
DROP TABLE IF EXISTS terms;
DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS folders;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
