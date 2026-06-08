-- +goose Up
-- Keep the newest proposal per client; soft-delete older duplicates.
UPDATE proposals AS p
SET deleted_at = NOW(),
    updated_at = NOW()
FROM (
    SELECT id,
        ROW_NUMBER() OVER (
            PARTITION BY client_id
            ORDER BY created_at DESC NULLS LAST, id DESC
        ) AS rn
    FROM proposals
    WHERE client_id IS NOT NULL
      AND deleted_at IS NULL
) AS ranked
WHERE p.id = ranked.id
  AND ranked.rn > 1;

DROP INDEX IF EXISTS idx_proposals_client_id;

CREATE UNIQUE INDEX IF NOT EXISTS uix_proposals_client_id
    ON proposals (client_id)
    WHERE deleted_at IS NULL
      AND client_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS uix_proposals_client_id;

CREATE INDEX IF NOT EXISTS idx_proposals_client_id ON proposals (client_id);
