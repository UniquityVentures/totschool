-- +goose Up
CREATE TABLE IF NOT EXISTS proposals (
    id                BIGSERIAL PRIMARY KEY,
    created_at        TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    created_by_id     BIGINT NOT NULL REFERENCES users (id),
    title             VARCHAR(250) NOT NULL,
    answers           JSONB,
    generated_content TEXT,
    generation_id     BIGINT
);

CREATE INDEX IF NOT EXISTS idx_proposals_deleted_at ON proposals (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS proposals;
