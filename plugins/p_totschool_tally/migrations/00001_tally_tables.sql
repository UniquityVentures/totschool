-- +goose Up
CREATE TABLE IF NOT EXISTS tot_school_sessions (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name   VARCHAR(250) UNIQUE,
    start  DATE,
    "end"  DATE
);

CREATE INDEX IF NOT EXISTS idx_tot_school_sessions_deleted_at ON tot_school_sessions (deleted_at);

CREATE TABLE IF NOT EXISTS tallies (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,
    user_id         BIGINT NOT NULL REFERENCES users (id),
    date            DATE NOT NULL,
    visits          INTEGER DEFAULT 0,
    appointments    INTEGER DEFAULT 0,
    leads           INTEGER DEFAULT 0,
    presentations   INTEGER DEFAULT 0,
    demos           INTEGER DEFAULT 0,
    letters         INTEGER DEFAULT 0,
    follow_ups      INTEGER DEFAULT 0,
    proposals       INTEGER DEFAULT 0,
    policies        INTEGER DEFAULT 0,
    premium         INTEGER DEFAULT 0,
    CONSTRAINT uix_tally_user_date UNIQUE (user_id, date)
);

CREATE INDEX IF NOT EXISTS idx_tallies_deleted_at ON tallies (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS tallies;
DROP TABLE IF EXISTS tot_school_sessions;
