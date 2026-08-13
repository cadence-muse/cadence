CREATE TABLE band
(
    id          UUID PRIMARY KEY,
    name        TEXT        NOT NULL,
    invite_code TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by  UUID        NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by  UUID
);
