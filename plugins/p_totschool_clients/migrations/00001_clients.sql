-- +goose Up
CREATE TABLE IF NOT EXISTS clients (
    id            BIGSERIAL PRIMARY KEY,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,
    created_by_id BIGINT NOT NULL REFERENCES users (id),
    name          VARCHAR(250) NOT NULL,
    address       TEXT,
    phone         VARCHAR(20),
    remarks       TEXT
);

CREATE INDEX IF NOT EXISTS idx_clients_deleted_at ON clients (deleted_at);
CREATE INDEX IF NOT EXISTS idx_clients_created_by_id ON clients (created_by_id);

-- +goose Down
DROP TABLE IF EXISTS clients;
