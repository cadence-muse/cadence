CREATE INDEX track_band_id_idx ON track (band_id) WHERE deleted_at IS NULL;
