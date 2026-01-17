CREATE TABLE albums (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  cover_image_url TEXT,
  release_date DATE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);