-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    CREATE TYPE appointment_status AS ENUM ('pending', 'done', 'cancelled', 'postponed');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END;
$$;
-- +goose StatementEnd

ALTER TABLE appointments
    ADD COLUMN IF NOT EXISTS client_id BIGINT REFERENCES clients (id);

ALTER TABLE appointments
    ADD COLUMN IF NOT EXISTS status appointment_status;

-- +goose StatementBegin
DO $$
DECLARE
    r RECORD;
    new_client_id BIGINT;
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = CURRENT_SCHEMA()
          AND table_name = 'appointments'
          AND column_name = 'name'
    ) THEN
        RETURN;
    END IF;

    FOR r IN
        SELECT id, created_at, updated_at, created_by_id, name, location, phone, remarks
        FROM appointments
        WHERE client_id IS NULL
        ORDER BY id
    LOOP
        INSERT INTO clients (created_at, updated_at, created_by_id, name, address, phone, remarks)
        VALUES (
            COALESCE(r.created_at, NOW()),
            COALESCE(r.updated_at, NOW()),
            r.created_by_id,
            r.name,
            NULLIF(BTRIM(r.location), ''),
            NULLIF(BTRIM(r.phone), ''),
            NULLIF(BTRIM(r.remarks), '')
        )
        RETURNING id INTO new_client_id;

        UPDATE appointments
        SET client_id = new_client_id
        WHERE id = r.id;
    END LOOP;
END;
$$;
-- +goose StatementEnd

UPDATE appointments
SET status = CASE
    WHEN datetime > NOW() THEN 'pending'::appointment_status
    WHEN datetime < NOW() THEN 'done'::appointment_status
    ELSE 'pending'::appointment_status
END
WHERE status IS NULL;

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM appointments WHERE client_id IS NULL
    ) THEN
        RAISE EXCEPTION 'appointments.client_id backfill incomplete: rows still NULL';
    END IF;
END;
$$;
-- +goose StatementEnd

ALTER TABLE appointments
    ALTER COLUMN client_id SET NOT NULL;

ALTER TABLE appointments
    ALTER COLUMN status SET NOT NULL;

ALTER TABLE appointments
    DROP COLUMN IF EXISTS name,
    DROP COLUMN IF EXISTS location,
    DROP COLUMN IF EXISTS phone;

CREATE INDEX IF NOT EXISTS idx_appointments_client_id ON appointments (client_id);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments (status);

-- +goose Down
ALTER TABLE appointments
    ADD COLUMN IF NOT EXISTS name VARCHAR(250),
    ADD COLUMN IF NOT EXISTS location TEXT,
    ADD COLUMN IF NOT EXISTS phone VARCHAR(20);

UPDATE appointments a
SET
    name = c.name,
    location = c.address,
    phone = c.phone
FROM clients c
WHERE c.id = a.client_id
  AND a.name IS NULL;

ALTER TABLE appointments
    ALTER COLUMN name SET NOT NULL;

ALTER TABLE appointments
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS client_id;

DROP TYPE IF EXISTS appointment_status;

DROP INDEX IF EXISTS idx_appointments_client_id;
DROP INDEX IF EXISTS idx_appointments_status;
