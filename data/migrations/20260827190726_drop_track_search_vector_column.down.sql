ALTER TABLE track ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (setweight(to_tsvector('simple', title), 'A') || setweight(to_tsvector('simple', artist), 'B')) STORED;
