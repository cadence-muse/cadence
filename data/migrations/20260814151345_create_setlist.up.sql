CREATE TABLE setlist
(
    id             UUID PRIMARY KEY,
    band_id        UUID        NOT NULL REFERENCES band (id) ON DELETE CASCADE,
    name           TEXT        NOT NULL,
    event_location TEXT,
    event_date     DATE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by     UUID        NOT NULL,
    updated_at     TIMESTAMPTZ,
    updated_by     UUID,
    deleted_at     TIMESTAMPTZ,
    deleted_by     UUID
);
