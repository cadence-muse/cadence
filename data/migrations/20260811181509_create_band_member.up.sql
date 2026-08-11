CREATE TABLE band_member
(
    band_id   UUID        NOT NULL REFERENCES band (id) ON DELETE CASCADE,
    user_id   UUID        NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    role      TEXT        NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (band_id, user_id)
);
