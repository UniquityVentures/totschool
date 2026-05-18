-- +goose Up
ALTER TABLE proposals
    ADD COLUMN IF NOT EXISTS client_id BIGINT REFERENCES clients (id);

CREATE INDEX IF NOT EXISTS idx_proposals_client_id ON proposals (client_id);

-- +goose Down
DROP INDEX IF EXISTS idx_proposals_client_id;

ALTER TABLE proposals
    DROP COLUMN IF EXISTS client_id;
