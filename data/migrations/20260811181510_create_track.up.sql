CREATE TABLE track
(
    id               UUID PRIMARY KEY,
    band_id          UUID        NOT NULL REFERENCES band (id) ON DELETE CASCADE,
    title            TEXT        NOT NULL,
    artist           TEXT        NOT NULL,
    duration_seconds INTEGER,
    tempo            SMALLINT,
    key              TEXT,
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by       UUID        NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by       UUID
);
