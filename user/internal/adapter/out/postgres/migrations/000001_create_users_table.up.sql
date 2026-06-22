CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT users_username_matches_email CHECK (username = email)
);