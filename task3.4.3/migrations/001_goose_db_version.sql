-- +goose Up
CREATE TABLE IF NOT EXISTS goose_db_version (
id BIGSERIAL PRIMARY KEY,
version_id VARCHAR(255) NOT NULL,
is_applied BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE IF EXISTS goose_db_version;
