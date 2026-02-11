CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id BIGINT REFERENCES posts(id),
    user_id BIGINT REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
