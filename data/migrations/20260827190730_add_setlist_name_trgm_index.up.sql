CREATE INDEX setlist_name_trgm_idx ON setlist USING GIN (name gin_trgm_ops) WHERE deleted_at IS NULL;
