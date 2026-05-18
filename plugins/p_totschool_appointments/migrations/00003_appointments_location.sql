-- +goose Up
ALTER TABLE appointments
    ADD COLUMN IF NOT EXISTS location TEXT;

-- +goose Down
ALTER TABLE appointments
    DROP COLUMN IF EXISTS location;
