CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    is_group BOOLEAN NOT NULL DEFAULT 0,
    title TEXT,
    photo_url TEXT,
    created_at DATETIME
);