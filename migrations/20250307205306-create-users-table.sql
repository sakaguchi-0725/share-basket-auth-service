-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    cognito_uid VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- +migrate Down
DROP TABLE IF EXISTS users;
