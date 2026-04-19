DROP TABLE IF EXISTS users;

CREATE TYPE role AS ENUM ('admin', 'user', 'guest');

CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name  VARCHAR(255) NOT NULL,
    role       role NOT NULL DEFAULT 'user',
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);
