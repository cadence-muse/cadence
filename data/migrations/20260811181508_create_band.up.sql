CREATE TABLE band
(
    id         UUID PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID        NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
