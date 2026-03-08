CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id SERIAL NOT NULL,
    sender_id TEXT NOT NULL,
    text TEXT,
    reply_to_id TEXT,
    created_at DATETIME,
    edited_at DATETIME,
    deleted_at DATETIME,
    read_at DATETIME
);