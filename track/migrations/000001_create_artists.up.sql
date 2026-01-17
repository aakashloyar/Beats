CREATE Table artists {
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    bio TEXT,
    profile_image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
};