CREATE TABLE likes (
    user_id BIGINT REFERENCES users(id),
    post_id BIGINT REFERENCES posts(id),
    PRIMARY KEY (user_id, post_id)
);
