CREATE TABLE track
(
    id               UUID PRIMARY KEY,
    band_id          UUID        NOT NULL REFERENCES band (id) ON DELETE CASCADE,
    title            TEXT        NOT NULL,
    artist           TEXT        NOT NULL,
    duration_seconds INTEGER     NOT NULL,
    original_tempo   SMALLINT    NOT NULL,
    original_key     TEXT        NOT NULL,
    custom_tempo     SMALLINT,
    custom_key       TEXT,
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by       UUID        NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
