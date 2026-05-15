-- +goose Up
CREATE TABLE IF NOT EXISTS appointments (
    id                BIGSERIAL PRIMARY KEY,
    created_at        TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    created_by_id     BIGINT NOT NULL REFERENCES users (id),
    name              VARCHAR(250) NOT NULL,
    location          TEXT,
    datetime          TIMESTAMPTZ NOT NULL,
    phone             VARCHAR(20),
    remarks           TEXT,
    extra_info        TEXT,
    generated_letter  TEXT,
    generation_id     BIGINT
);

CREATE INDEX IF NOT EXISTS idx_appointments_deleted_at ON appointments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_appointments_created_by_id ON appointments (created_by_id);
CREATE INDEX IF NOT EXISTS idx_appointments_datetime ON appointments (datetime);

-- +goose Down
DROP TABLE IF EXISTS appointments;
