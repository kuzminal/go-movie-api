CREATE EXTENSION pg_trgm;
CREATE INDEX IF NOT EXISTS movies_title_idx ON movies USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS movies_genres_idx ON movies USING GIN (genres);