CREATE TABLE band
(
    id          UUID PRIMARY KEY,
    name        TEXT        NOT NULL,
    invite_code TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by  UUID        NOT NULL,
    updated_at  TIMESTAMPTZ,
    updated_by  UUID,
    deleted_at  TIMESTAMPTZ,
    deleted_by  UUID
);
