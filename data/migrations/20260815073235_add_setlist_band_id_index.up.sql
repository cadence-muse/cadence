CREATE INDEX setlist_band_id_idx ON setlist (band_id) WHERE deleted_at IS NULL;
