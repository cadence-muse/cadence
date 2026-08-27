CREATE INDEX track_title_artist_trgm_idx ON track USING GIN ((title || ' ' || artist) gin_trgm_ops) WHERE deleted_at IS NULL;
