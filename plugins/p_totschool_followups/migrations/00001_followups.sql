-- +goose Up
CREATE TABLE IF NOT EXISTS followups (
    id                BIGSERIAL PRIMARY KEY,
    created_at        TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    created_by_id     BIGINT NOT NULL REFERENCES users (id),
    client_id         BIGINT NOT NULL REFERENCES clients (id),
    title             VARCHAR(250) NOT NULL,
    extra_info        TEXT,
    generated_letter  TEXT,
    generation_id     BIGINT
);

CREATE INDEX IF NOT EXISTS idx_followups_deleted_at ON followups (deleted_at);
CREATE INDEX IF NOT EXISTS idx_followups_created_by_id ON followups (created_by_id);
CREATE INDEX IF NOT EXISTS idx_followups_client_id ON followups (client_id);

-- +goose Down
DROP TABLE IF EXISTS followups;
