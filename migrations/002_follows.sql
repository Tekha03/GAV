CREATE TABLE follows (
    follower_id BIGINT REFERENCES users(id),
    following_id BIGINT REFERENCES users(id),
    PRIMARY KEY (follower_id, following_id)
);
