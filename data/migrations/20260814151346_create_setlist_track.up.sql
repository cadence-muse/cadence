CREATE TABLE setlist_track
(
    setlist_id UUID    NOT NULL REFERENCES setlist (id) ON DELETE CASCADE,
    track_id   UUID    NOT NULL REFERENCES track (id) ON DELETE CASCADE,
    position   INTEGER NOT NULL,
    PRIMARY KEY (setlist_id, track_id),
    UNIQUE (setlist_id, position) DEFERRABLE INITIALLY DEFERRED
);
