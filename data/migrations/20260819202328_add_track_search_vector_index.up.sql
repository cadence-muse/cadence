CREATE INDEX track_search_vector_idx ON track USING GIN (search_vector) WHERE deleted_at IS NULL;
