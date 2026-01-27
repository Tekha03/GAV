CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

ALTER TABLE users ADD COLUMN role_id INT REFERENCES roles(id);

INSERT INTO roles (name) VALUES ('user'), ('admin');
