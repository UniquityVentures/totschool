-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    CREATE TYPE client_status AS ENUM ('active', 'archived');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END;
$$;
-- +goose StatementEnd

ALTER TABLE clients
    ADD COLUMN IF NOT EXISTS status client_status;

UPDATE clients
SET status = 'active'::client_status
WHERE status IS NULL;

ALTER TABLE clients
    ALTER COLUMN status SET DEFAULT 'active'::client_status,
    ALTER COLUMN status SET NOT NULL;

-- +goose Down
ALTER TABLE clients DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS client_status;
