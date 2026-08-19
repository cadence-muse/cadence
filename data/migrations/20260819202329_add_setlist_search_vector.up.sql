ALTER TABLE setlist ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (to_tsvector('simple', name)) STORED;
