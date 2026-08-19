CREATE INDEX setlist_search_vector_idx ON setlist USING GIN (search_vector) WHERE deleted_at IS NULL;
