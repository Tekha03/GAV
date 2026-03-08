CREATE TABLE chat_members (
    user_id SERIAL NOT NULL,
    chat_id SERIAL NOT NULL,
    role TEXT NOT NULL,
    joined_at DATETIME,
    left_at DATETIME,
    muted BOOLEAN DEFAULT 0,
    last_read_message_id TEXT,
    PRIMARY KEY (user_id, chat_id)
);